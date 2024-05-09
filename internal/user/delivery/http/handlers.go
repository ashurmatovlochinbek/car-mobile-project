package http

import (
	"car-mobile-project/config"
	"car-mobile-project/internal/models"
	"car-mobile-project/internal/user"
	"car-mobile-project/pkg/util"
	"encoding/json"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

type userHandler struct {
	cfg    *config.Config
	userUC user.UseCase
}

func NewUserHandler(cfg *config.Config, userUC user.UseCase) user.Handlers {
	return &userHandler{
		cfg:    cfg,
		userUC: userUC,
	}
}

func (uh *userHandler) VerifyOTP() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var otp *models.OTPObject
		var responseObject *models.ResponseObject
		err := json.NewDecoder(r.Body).Decode(&otp)

		if err != nil {
			responseObject = models.NewResponseObject(
				false,
				err.Error(),
				nil,
			)
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(responseObject)
			w.Write(jsonResponse)
			return
		}

		if otp.OTPPassword == "4444" {

			u, err := uh.userUC.GetByPhoneNumber(r.Context(), otp.PhoneNumber)

			if err != nil {
				responseObject = models.NewResponseObject(false, err.Error(), nil)
				jsonResponse, _ := json.Marshal(responseObject)
				re := regexp.MustCompile(`status:\s*(\d+)`)
				match := re.FindStringSubmatch(err.Error())
				statusCode, _ := strconv.Atoi(match[1])
				w.WriteHeader(statusCode)
				w.Write(jsonResponse)
				return
			}

			accessTokenString, err := util.GenerateJWTToken(u, uh.cfg)

			if err != nil {
				responseObject = models.NewResponseObject(false, err.Error(), nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonResponse)
				return
			}

			refreshTokenString, err := util.GenerateJWTToken(u, uh.cfg)

			if err != nil {
				responseObject = models.NewResponseObject(false, err.Error(), nil)
				jsonResponse, _ := json.Marshal(responseObject)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write(jsonResponse)
				return
			}

			responseObject = models.NewResponseObject(true, "Access and Refresh tokens are valid", map[string]string{
				"access_token":  accessTokenString,
				"refresh_token": refreshTokenString,
			})
			jsonResponse, _ := json.Marshal(responseObject)
			w.Write(jsonResponse)
		} else {
			responseObject = models.NewResponseObject(false, "otp is incorrect", nil)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusBadRequest)
			w.Write(jsonResponse)
			return
		}
	}
}

func (uh *userHandler) Refresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		splitToken := strings.Split(tokenString, "Bearer ")
		reqToken := splitToken[1]
		token, _ := jwt.Parse(reqToken, func(token *jwt.Token) (interface{}, error) {
			return []byte(uh.cfg.Server.JwtSecretKey), nil
		})

		tokenClaims := token.Claims.(jwt.MapClaims)

		userId, _ := uuid.Parse(tokenClaims["user_id"].(string))
		name := tokenClaims["name"].(string)
		phoneNumber := tokenClaims["phone_number"].(string)

		accessTokenString, _ := util.GenerateJWTToken(&models.User{
			userId,
			name,
			phoneNumber,
		}, uh.cfg)

		refreshTokenString, _ := util.GenerateRefreshToken(&models.User{
			userId,
			name,
			phoneNumber,
		}, uh.cfg)

		responseObject := models.NewResponseObject(true, "access and refresh tokens", map[string]string{
			"access_token":  accessTokenString,
			"refresh_token": refreshTokenString,
		})
		jsonResponse, _ := json.Marshal(responseObject)
		w.Write(jsonResponse)
	}
}

func (uh *userHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var responseObject *models.ResponseObject
		var user models.User

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			responseObject = models.NewResponseObject(
				false,
				err.Error(),
				nil,
			)
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(responseObject)
			w.Write(jsonResponse)
			return
		}

		u, err := uh.userUC.GetByPhoneNumber(r.Context(), user.PhoneNumber)

		if err != nil {
			responseObject = models.NewResponseObject(false, err.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			re := regexp.MustCompile(`status:\s*(\d+)`)
			match := re.FindStringSubmatch(err.Error())
			statusCode, _ := strconv.Atoi(match[1])
			w.WriteHeader(statusCode)
			w.Write(jsonResponse)
			return
		} else {
			responseObject = models.NewResponseObject(true, "Login is success", u)
			jsonResponse, _ := json.Marshal(responseObject)
			w.WriteHeader(http.StatusOK)
			w.Write(jsonResponse)
		}

		//accessTokenString, err := util.GenerateJWTToken(u, uh.cfg)
		//
		//if err != nil {
		//	responseObject = models.NewResponseObject(false, err.Error(), nil)
		//	jsonResponse, _ := json.Marshal(responseObject)
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write(jsonResponse)
		//	return
		//}
		//
		//refreshTokenString, err := util.GenerateJWTToken(u, uh.cfg)
		//
		//if err != nil {
		//	responseObject = models.NewResponseObject(false, err.Error(), nil)
		//	jsonResponse, _ := json.Marshal(responseObject)
		//	w.WriteHeader(http.StatusInternalServerError)
		//	w.Write(jsonResponse)
		//	return
		//}

		//responseObject = models.NewResponseObject(true, "Access and Refresh tokens are valid", map[string]string{
		//	"access_token":  accessTokenString,
		//	"refresh_token": refreshTokenString,
		//})
		//jsonResponse, _ := json.Marshal(responseObject)
		//w.Write(jsonResponse)
	}
}

func (uh *userHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var responseObject *models.ResponseObject
		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			responseObject = models.NewResponseObject(
				false,
				err.Error(),
				nil,
			)
			w.WriteHeader(http.StatusBadRequest)
			jsonResponse, _ := json.Marshal(responseObject)
			w.Write(jsonResponse)
			return
		}

		createdUser, err := uh.userUC.Create(r.Context(), user)

		if err != nil {
			responseObject = models.NewResponseObject(false, err.Error(), nil)
			jsonResponse, _ := json.Marshal(responseObject)
			re := regexp.MustCompile(`status:\s*(\d+)`)
			match := re.FindStringSubmatch(err.Error())
			statusCode, _ := strconv.Atoi(match[1])
			w.WriteHeader(statusCode)
			w.Write(jsonResponse)
			return
		}

		w.WriteHeader(http.StatusCreated)
		responseObject = models.NewResponseObject(true, "user created", createdUser)
		jsonResponse, _ := json.Marshal(responseObject)
		w.Write(jsonResponse)
	}
}

func (uh *userHandler) GetSecuredResource() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello"))
	}
}
