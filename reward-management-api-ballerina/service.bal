import ballerina/http;
import ballerina/log;
import ballerina/os;

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

public type User record {
    string userId;
    string firstName;
    string lastName;
    string email;
};

public type Reward record {
    string rewardId;
    string userId;
    string firstName;
    string lastName;
    string email;
};



configurable string vendorManagementClientId = os:getEnv("VENDOR_MANAGEMENT_CLIENT_ID"); 
configurable string vendorManagementClientSecret = os:getEnv("VENDOR_MANAGEMENT_CLIENT_SECRET"); 
configurable string vendorManagementTokenUrl = os:getEnv("VENDOR_MANAGEMENT_TOKEN_URL"); 
configurable string vendorManagementApiUrl = os:getEnv("VENDOR_MANAGEMENT_API_URL"); 

configurable string loyaltyClientId = os:getEnv("LOYALTY_CLIENT_ID"); 
configurable string loyaltyClientSecret = os:getEnv("LOYALTY_CLIENT_SECRET");
configurable string loyaltyTokenUrl = os:getEnv("LOYALTY_TOKEN_URL"); 
configurable string loyaltyApiUrl = os:getEnv("LOYALTY_API_URL");


# The client to connect to the loyalty management api
@display {
    label: "loyalty management api",
    id: "loyalty-management-api"
}
http:Client loyaltyAPIEndpoint = check new (loyaltyApiUrl, {
    auth: {
        tokenUrl: loyaltyTokenUrl,
        clientId: loyaltyClientId,
        clientSecret: loyaltyClientSecret
    }
});

# The client to connect to the vendor management api
@display {
    label: "vendor management api",
    id: "vendor-management-api"
}
http:Client vendorManagementClientEp = check new (vendorManagementApiUrl, {
    auth: {
        tokenUrl: vendorManagementTokenUrl,
        clientId: vendorManagementClientId,
        clientSecret: vendorManagementClientSecret
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
        if (user is http:Error) {
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

    firstName: user.firstName,
    lastName: user.lastName,
    userId: rewardSelection.userId,
    rewardId: rewardSelection.selectedRewardDealId,
    email: user.email
};
