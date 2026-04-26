package auth

import (
"crypto/rand"
"encoding/hex"
"errors"
"time"

"github.com/gin-gonic/gin"
"github.com/golang-jwt/jwt/v5"
"golang.org/x/crypto/bcrypt"
"gitlab.franciscan.edu/spr-26-backend/basement/internal/database/queries"
"gitlab.franciscan.edu/spr-26-backend/basement/internal/models"
)

var jwtSecret = []byte("your-secret-key") // In production, use env var

type Claims struct {
UserID string `json:"user_id"`
jwt.RegisteredClaims
}

type LoginRequest struct {
Email    string `json:"email" binding:"required"`
Password string `json:"password" binding:"required"`
}

type SignupRequest struct {
Username string `json:"username" binding:"required"`
Email    string `json:"email" binding:"required"`
Password string `json:"password" binding:"required"`
DOB      string `json:"dob" binding:"required"`
}

func HashPassword(password string) (string, error) {
bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
return err == nil
}

func GenerateJWT(userID string) (string, error) {
claims := Claims{
UserID: userID,
RegisteredClaims: jwt.RegisteredClaims{
ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
IssuedAt:  jwt.NewNumericDate(time.Now()),
},
}
token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
return token.SignedString(jwtSecret)
}

func ValidateJWT(tokenString string) (*Claims, error) {
token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
return jwtSecret, nil
})
if err != nil {
return nil, err
}
if claims, ok := token.Claims.(*Claims); ok && token.Valid {
return claims, nil
}
return nil, errors.New("invalid token")
}

func AuthMiddleware() gin.HandlerFunc {
return func(c *gin.Context) {
tokenString := c.GetHeader("Authorization")
if tokenString == "" {
c.JSON(401, gin.H{"error": "Authorization header required"})
c.Abort()
return
}
// Remove "Bearer " prefix if present
if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
tokenString = tokenString[7:]
}
claims, err := ValidateJWT(tokenString)
if err != nil {
c.JSON(401, gin.H{"error": "Invalid token"})
c.Abort()
return
}
c.Set("user_id", claims.UserID)
c.Next()
}
}

func Signup(c *gin.Context) {
var req SignupRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(400, gin.H{"error": err.Error()})
return
}

hashedPassword, err := HashPassword(req.Password)
if err != nil {
c.JSON(500, gin.H{"error": "Failed to hash password"})
return
}

userID := generateID()
user := models.UserModel{
ID:           userID,
Name:         req.Username,
Email:        req.Email,
PasswordHash: hashedPassword,
DOB:          req.DOB,
}

if err := queries.CreateUser(user); err != nil {
c.JSON(500, gin.H{"error": "Failed to create user"})
return
}

token, err := GenerateJWT(userID)
if err != nil {
c.JSON(500, gin.H{"error": "Failed to generate token"})
return
}

c.JSON(201, gin.H{"token": token, "user": user})
}

func Login(c *gin.Context) {
var req LoginRequest
if err := c.ShouldBindJSON(&req); err != nil {
c.JSON(400, gin.H{"error": err.Error()})
return
}

user, err := queries.GetUserByEmail(req.Email)
if err != nil || !CheckPasswordHash(req.Password, user.PasswordHash) {
c.JSON(401, gin.H{"error": "Invalid credentials"})
return
}

token, err := GenerateJWT(user.ID)
if err != nil {
c.JSON(500, gin.H{"error": "Failed to generate token"})
return
}

c.JSON(200, gin.H{"token": token})
}

func generateID() string {
bytes := make([]byte, 16)
rand.Read(bytes)
return hex.EncodeToString(bytes)
}
