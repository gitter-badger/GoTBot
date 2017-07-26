/**
 *
 * This file is used to provide type hinting in your javascript plugins.
 * It is as well a documentation which functions/global vars you can use.
 * It is not used in any other way and is never called.
 *
 */


/**
 * The name of the user who triggered the event
 */
var sender;

/**
 * Sends a message to the channel where the current event was triggered in
 * @param {string} message
 */
function sendMessage(message) {
}

/**
 * Returns a user data object as json string.
 * Use JSON.parse() to get a json object.
 * @param {string} name
 */
function getUser(name) {
    return {
        Name: "",
        MessageCount: 0,
        LastJoin: "",
        LastPart: "",
        LastActive: "",
        FirstSeen: ""
    }
}