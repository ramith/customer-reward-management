package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"golang.org/x/oauth2/clientcredentials"
	"net/http"
	"os"
)

type User struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type RewardSelection struct {
	UserId               string `json:"userId"`
	SelectedRewardDealId string `json:"selectedRewardDealId"`
	AcceptedTnC          bool   `json:"acceptedTnC"`
}

type Reward struct {
	RewardId  string `json:"rewardId"`
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

var logger *zap.Logger

var vendorManagementClientId = os.Getenv("VENDOR_MANAGEMENT_CLIENT_ID")
var vendorManagementClientSecret = os.Getenv("VENDOR_MANAGEMENT_CLIENT_SECRET")
var vendorManagementTokenUrl = os.Getenv("VENDOR_MANAGEMENT_TOKEN_URL")
var vendorManagementApiUrl = os.Getenv("VENDOR_MANAGEMENT_API_URL")

var loyaltyClientId = os.Getenv("LOYALTY_CLIENT_ID")
var loyaltyClientSecret = os.Getenv("LOYALTY_CLIENT_SECRET")
var loyaltyTokenUrl = os.Getenv("LOYALTY_TOKEN_URL")
var loyaltyApiUrl = os.Getenv("LOYALTY_API_URL")

var vendorManagementClientCredsConfig = clientcredentials.Config{
	ClientID:     vendorManagementClientId,
	ClientSecret: vendorManagementClientSecret,
	TokenURL:     vendorManagementTokenUrl,
}

var loyaltyManagementClientCredsConfig = clientcredentials.Config{
	ClientID:     loyaltyClientId,
	ClientSecret: loyaltyClientSecret,
	TokenURL:     loyaltyTokenUrl,
}


func HandleRewardSelection(w http.ResponseWriter, r *http.Request) {
	var selection RewardSelection

	// Decode the request body into the RewardSelection struct
	if err := json.NewDecoder(r.Body).Decode(&selection); err != nil {
		logger.Error("Failed to decode incoming reward selection data", zap.Error(err))
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid request payload"))
		return
	}

	logger.Info("received reward selection",
		zap.String("UserId", selection.UserId),
		zap.String("SelectedRewardDealId", selection.SelectedRewardDealId),
		zap.Bool("AcceptedTnC", selection.AcceptedTnC),
	)

	user, err := FetchUserByIdFromLoyaltyApi(selection.UserId)
	if err != nil {
		logger.Error("Failed to fetch user", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("failed to fetch user"))
		return
	}

	err = PostRewardSelectionToVendorManagementApi(Reward{
		RewardId:  selection.SelectedRewardDealId,
		UserId:    selection.UserId,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
	})

	if err != nil {
		logger.Error("unable to send reward selection to vendor management api", zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to send reward selection to vendor management  api"))
		return
	}

	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("reward selection received successfully"))
}

func LivenessProbe(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Alive"))
}

func ReadinessProbe(w http.ResponseWriter, r *http.Request) {
	// Add logic here to check database connections, external services, etc.
	// If all checks pass:
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Ready"))
	// If any check fails:
	// w.WriteHeader(http.StatusInternalServerError)
}

func main() {

	defer logger.Sync() // Ensure all buffered logs are written

	logger.Info("starting the reward management api (golang)...")

	r := mux.NewRouter()
	r.HandleFunc("/select-reward", HandleRewardSelection).Methods("POST")

	r.HandleFunc("/healthz", ReadinessProbe).Methods("GET") // Readiness probe
	r.HandleFunc("/livez", LivenessProbe).Methods("GET")    // Liveness probe

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		
		return
	}
}

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func FetchUserByIdFromLoyaltyApi(userId string) (*User, error) {
	// Construct the full URL using the base URL from the environment variable
	url := fmt.Sprintf("%s/user/%s", loyaltyApiUrl, userId)
	// Make the HTTP GET request
	resp, err := loyaltyManagementClientCredsConfig.Client(context.Background()).Get(url)
	if err != nil {
		logger.Error("Failed to fetch user", zap.String("userId", userId), zap.Error(err))
		return nil, fmt.Errorf("failed to fetch user: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		logger.Warn("API responded with non-200 status code", zap.Int("statusCode", resp.StatusCode))
		return nil, fmt.Errorf("API responded with status code: %d", resp.StatusCode)
	}

	// Decode the response body into the User struct
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		logger.Error("Failed to decode user data", zap.Error(err))
		return nil, fmt.Errorf("failed to decode user data: %v", err)
	}

	logger.Info("Successfully fetched user", zap.String("userId", user.UserId))
	return &user, nil
}

func PostRewardSelectionToVendorManagementApi(reward Reward) error {
	// Marshal the Reward struct to JSON
	payload, err := json.Marshal(reward)
	if err != nil {
		logger.Error("Failed to marshal reward", zap.Error(err))
		return err
	}

	// Construct the full URL using the base URL from the environment variable
	url := fmt.Sprintf("%s/rewards", vendorManagementApiUrl)
	// Make the HTTP POST request
	resp, err := vendorManagementClientCredsConfig.Client(context.Background()).Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		logger.Error("Failed to send reward to vendor management", zap.Any("reward", payload), zap.Error(err))
		return fmt.Errorf("Failed to send reward to vendor management: %v", err)
	}
	defer resp.Body.Close()

	// Check for non-200 status codes
	if resp.StatusCode != http.StatusOK {
		logger.Warn("API responded with non-200 status code", zap.Int("statusCode", resp.StatusCode))
		return fmt.Errorf("API responded with status code: %d", resp.StatusCode)
	}

	logger.Info("Successfully sent reward selection user", zap.Any("reward", reward))
	return nil
}
