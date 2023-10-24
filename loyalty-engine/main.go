/*
 * Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
 *
 * WSO2 LLC. licenses this file to you under the Apache License,
 * Version 2.0 (the "License"); you may not use this file except
 * in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type User struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email"`
}

type RewardOffer struct {
	Id          string  `json:"id"`
	Name        string  `json:"name"`
	Value       float32 `json:"value"`
	TotalPoints int     `json:"totalPoints"`
	Description string  `json:"description"`
	LogoUrl     string  `json:"logoUrl"`
}

type UserReward struct {
	UserId               string `json:"userId"`
	SelectedRewardDealId string `json:"selectedRewardDealId"`
	Timestamp            string `json:"timestamp"`
	AcceptedTnC          bool   `json:"acceptedTnC"`
}

var logger *zap.Logger
var userRewards []UserReward
var rewardOffers []RewardOffer
var users []User

func init() {
	var err error
	logger, err = zap.NewProduction()
	if err != nil {
		panic(err)
	}
}

func main() {

	defer logger.Sync() // Ensure all buffered logs are written

	logger.Info("Starting the loyalty engine...")

	r := mux.NewRouter()

	rewardOffers = append(rewardOffers, RewardOffer{"RWD34589", "Target", 25, 500, "A Target GiftCard is your opportunity to shop for thousands of items at more than 1,900 Target stores in the U.S., as well as Target.com. From home décor, small appliances and electronics to fashion, accessories and music, find exactly what you’re looking for at Target. No fees. No expiration. No kidding.™", "https://drive.google.com/file/d/1FEOGLEG99HsttPBXliXUi8aWYqnNmPH2/view?usp=drive_link"})
	rewardOffers = append(rewardOffers, RewardOffer{"RWD34590", "Starbucks Coffee", 15, 200, "Enjoy a PM pick-me-up with a lunch sandwich, protein box or a bag of coffee—including Starbucks VIA Instant", "https://drive.google.com/file/d/1nku2n63zXBfrA3Bf0eAVWqu45mFLnaRE/view?usp=drive_link"})
	rewardOffers = append(rewardOffers, RewardOffer{"RWD34591", "Jumba Juice", 6, 600, "Let Jamba come to you – wherever you are. Get our Whirld Famous smoothies, juices, and bowls delivered in just a few clicks. My Jamba rewards members can also apply rewards & earn points on delivery orders when you order on jamba.com or the jamba app!", "https://drive.google.com/file/d/1khJX-N7N8xHrV5o9GvqsH7wApDoY8ej0/view?usp=drive_link"})
	rewardOffers = append(rewardOffers, RewardOffer{"RWD34592", "Grubhub", 10, 500, "Grubhub offers quick, easy food delivery, either online or through a mobile app. Customers can select from any local participating restaurant. They can add whatever they like to their order and have it delivered right to their home or office by one of Grubhub's delivery drivers. You can save even more by using a Grubhub promo code on your order", "https://drive.google.com/file/d/14S6olzLfOQJatEr4FkXyB_m1l31H2XyJ/view?usp=drive_link"})

	userRewards = append(userRewards, UserReward{"U451298", "RWD34589", "2023-09-04T14:32:21Z", true})
	userRewards = append(userRewards, UserReward{"U451299", "RWD34590", "2023-09-04T14:32:21Z", true})

	users = append(users, User{"U451298", "John", "Doe", "john@example.com"})
	users = append(users, User{"U451299", "Katie", "Smith", "katie@example.com"})
	users = append(users, User{"U451300", "Peter", "Parker", "peter@example.com"})

	r.HandleFunc("/rewards", getRewardOffers).Methods("GET")
	r.HandleFunc("/rewards/{id}", getRewardOffer).Methods("GET")
	r.HandleFunc("/user-rewards", getUserRewards).Methods("GET")
	r.HandleFunc("/user/{id}", getUserDetails).Methods("GET")
	http.ListenAndServe(":8080", r)
}

func getRewardOffers(w http.ResponseWriter, r *http.Request) {
	logger.Info("get reward offers")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rewardOffers)
}

func getRewardOffer(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range rewardOffers {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			logger.Info("get reward offer", zap.Any("reward offer", item))
			return
		}
	}

	logger.Info("reward offer not found", zap.String("offer id", params["id"]))
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&User{})
}

func getUserRewards(w http.ResponseWriter, r *http.Request) {
	logger.Info("get all rewards")
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userRewards)
}

func getUserDetails(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, user := range users {
		if user.UserId == params["id"] {
			json.NewEncoder(w).Encode(user)
			logger.Info("get user details", zap.Any("user", user))
			return
		}
	}

	logger.Info("user not found", zap.String("user id", params["id"]))
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&User{})
}
