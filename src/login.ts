import {readable} from "svelte/store";
import {PublicClientApplication, LogLevel, InteractionRequiredAuthError} from "@azure/msal-browser";

export let loggedInUser;

let loginChangeHook = (user) => {};
let hardAzureLogout = false;
function setLoginChangeHook(f) {
    loginChangeHook = f;
}

loggedInUser = readable(null, function(set) {
    setLoginChangeHook((user) => {
        console.dir(user);
        set(user);
    });

    return () => {}     // nothing to do when last subscriber unsubscribes
});

export function azureLogin() {
    signIn();
}

export function azureLogout() {
    signOut();
}

interface UserDef {
    username: string;
    name: string;
    token: string;
    idToken: string;
}

/**
 * Configuration object to be passed to MSAL instance on creation. 
 * For a full list of MSAL.js configuration parameters, visit:
 * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/configuration.md 
 */
 const msalConfig = {
    auth: {
        clientId: "a2e4c97f-dfe8-4945-9e6a-aa5973cea657",
        authority: "https://login.microsoftonline.com/db8e2ba9-95c1-4fbb-b558-6bf8bb1d2981",
        redirectUri: "http://localhost:5000",
    },
    cache: {
        cacheLocation: "sessionStorage", // This configures where your cache will be stored
        storeAuthStateInCookie: false, // Set this to "true" if you are having issues on IE11 or Edge
    },
    system: {	
        loggerOptions: {	
            loggerCallback: (level, message, containsPii) => {	
                if (containsPii) {		
                    return;		
                }		
                switch (level) {		
                    case LogLevel.Error:		
                        console.error(message);		
                        return;		
                    case LogLevel.Info:		
                        console.info(message);		
                        return;		
                    case LogLevel.Verbose:		
                        console.debug(message);		
                        return;		
                    case LogLevel.Warning:		
                        console.warn(message);		
                        return;		
                }	
            }	
        }	
    }
};

/**
 * Scopes you add here will be prompted for user consent during sign-in.
 * By default, MSAL.js will add OIDC scopes (openid, profile, email) to any login request.
 * For more information about OIDC scopes, visit: 
 * https://docs.microsoft.com/en-us/azure/active-directory/develop/v2-permissions-and-consent#openid-connect-scopes
 */
const loginRequest = {
    scopes: ["User.Read"]
};

/**
 * Add here the scopes to request when obtaining an access token for MS Graph API. For more information, see:
 * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/resources-and-scopes.md
 */
const tokenRequest = {
    scopes: ["User.Read", "Mail.Read"],
    forceRefresh: false // Set this to "true" to skip a cached token and go to the server to get a new token
};

// Create the main myMSALObj instance
// configuration parameters are located at authConfig.js
const myMSALObj = new PublicClientApplication(msalConfig);

let username = "";

function showWelcomeMessage(username) {
    alert('Welcome ' + username)
}

function selectAccount() {

    /**
     * See here for more info on account retrieval: 
     * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-common/docs/Accounts.md
     */

    const currentAccounts = myMSALObj.getAllAccounts();
    if (currentAccounts.length === 0) {
        return;
    } else if (currentAccounts.length > 1) {
        // Add choose account code here
        console.warn("Multiple accounts detected.");
        console.dir(currentAccounts);
    } else if (currentAccounts.length === 1) {
        username = currentAccounts[0].username;
        showWelcomeMessage(username);
        console.warn("select account is broken");
        console.dir(currentAccounts);
    }
}

function handleResponse(response) {

    /**
     * To see the full list of response object properties, visit:
     * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/request-response-object.md#response
     */

    if (response !== null) {
        console.info("logged in.");
        console.dir(response);
        let user: UserDef = {
            username: response.account.username,
            name: response.account.name,
            token: response.accessToken,
            idToken: response.idToken
        };
        loginChangeHook(user);
    } else {
        selectAccount();
    }
}

function signIn() {

    /**
     * You can pass a custom request object below. This will override the initial configuration. For more information, visit:
     * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/request-response-object.md#request
     */

    myMSALObj.loginPopup(loginRequest)
        .then(handleResponse)
        .catch(error => {
            console.error(error);
        });
}

function signOut() {

    /**
     * You can pass a custom request object below. This will override the initial configuration. For more information, visit:
     * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-browser/docs/request-response-object.md#request
     */

    const logoutRequest = {
        account: myMSALObj.getAccountByUsername(username),
        postLogoutRedirectUri: msalConfig.auth.redirectUri,
        mainWindowRedirectUri: msalConfig.auth.redirectUri
    };

    if (hardAzureLogout) {
        myMSALObj.logoutPopup(logoutRequest);   // this would log Azure Portal out.
    }
    loginChangeHook(null);
}

// function getTokenPopup(request) {

//     /**
//      * See here for more info on account retrieval: 
//      * https://github.com/AzureAD/microsoft-authentication-library-for-js/blob/dev/lib/msal-common/docs/Accounts.md
//      */
//     request.account = myMSALObj.getAccountByUsername(username);
    
//     return myMSALObj.acquireTokenSilent(request)
//         .catch(error => {
//             console.warn("silent token acquisition fails. acquiring token using popup");
//             if (error instanceof InteractionRequiredAuthError) {
//                 // fallback to interaction when silent call fails
//                 return myMSALObj.acquireTokenPopup(request)
//                     .then(tokenResponse => {
//                         console.log(tokenResponse);
//                         return tokenResponse;
//                     }).catch(error => {
//                         console.error(error);
//                     });
//             } else {
//                 console.warn(error);   
//             }
//     });
// }

// function seeProfile() {
//     getTokenPopup(loginRequest)
//         .then(response => {
//             callMSGraph(graphConfig.graphMeEndpoint, response.accessToken, updateUI);
//         }).catch(error => {
//             console.error(error);
//         });
// }

// function readMail() {
//     getTokenPopup(tokenRequest)
//         .then(response => {
//             callMSGraph(graphConfig.graphMailEndpoint, response.accessToken, updateUI);
//         }).catch(error => {
//             console.error(error);
//         });
// }

// selectAccount();
