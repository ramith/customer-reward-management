// Copyright (c) 2023, WSO2 LLC. (https://www.wso2.com/) All Rights Reserved.
//
// WSO2 LLC. licenses this file to you under the Apache License,
// Version 2.0 (the "License"); you may not use this file except
// in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing,
// software distributed under the License is distributed on an
// "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
// KIND, either express or implied. See the License for the
// specific language governing permissions and limitations
// under the License.
//

import ballerina/http;
import ballerina/log;
import ballerina/oauth2;

# Represents selected reward.
#
# + userId - id of the user 
# + selectedRewardDealId - id of the selected reward deal  
# + acceptedTnC - indicated weather the user accepted the terms and conditions
public type RewardSelection record {
    string userId;
    string selectedRewardDealId;
    boolean acceptedTnC;
};

# Represents a user.
#
# + userId - id of the user
# + firstName - first name of the user
# + lastName - last name of the user
# + email - email of the user
public type User record {
    string userId;
    string firstName;
    string lastName;
    string email;
};

# Represents a reward.
#
# + rewardId - id of the reward 
# + userId - id of the user 
# + firstName - first name of the user
# + lastName - last name of the user 
# + email - email of the user
public type Reward record {
    string rewardId;
    string userId;
    string firstName;
    string lastName;
    string email;
};

configurable string clientId = ?;
configurable string clientSecret = ?;
configurable string tokenUrl = ?;
configurable string loyaltyApiUrl = ?;
configurable string vendorManagementApiUrl = ?;

oauth2:ClientOAuth2Provider provider = new ({
    tokenUrl: tokenUrl,
    clientId: clientId,
    clientSecret: clientSecret
});

# The client to connect to the Loyalty Management API
@display {
    label: "Loyalty Management API",
    id: "loyalty-management-api"
}
http:Client loyaltyAPIEndpoint = check new (loyaltyApiUrl, {
    auth: {
        tokenUrl: tokenUrl,
        clientId: clientId,
        clientSecret: clientSecret
    }
});

# The client to connect to the Vendor Management API
@display {
    label: "Vendor Management API",
    id: "vendor-management-api"
}
http:Client vendorManagementClientEp = check new (vendorManagementApiUrl, {
    auth: {
        tokenUrl: tokenUrl,
        clientId: clientId,
        clientSecret: clientSecret
    }
});

# A service representing a network-accessible API
# bound to port `9090`.
@display {
    label: "Rewards Management API (Ballerina Implementation)",
    id: "rewards-management-api-ballerina"
}
service / on new http:Listener(9090) {

    resource function post select\-reward(RewardSelection selection) returns error|string {
        log:printInfo("reward selected: ", selection = selection);

        User|http:Error user = loyaltyAPIEndpoint->/user/[selection.userId];
        if user is http:Error {
            log:printError("error retrieving user: ", 'error = user);
            return user;
        }

        log:printInfo("user retrieved: ", user = user);
        Reward reward = transform(user, selection);

        http:Response|http:Error response = vendorManagementClientEp->post("/rewards", reward);
        if response is http:Error {
            log:printError("error while sending reward selection to vender ", 'error = response);
            return response;
        }

        log:printInfo("reward selection sent to vendor ", statusCode = response.statusCode);
        return "success";
    }

}

function transform(User user, RewardSelection rewardSelection) returns Reward => {
    rewardId: rewardSelection.selectedRewardDealId,
    userId: rewardSelection.userId,
    firstName: user.firstName,
    lastName: user.lastName,
    email: user.email
};
