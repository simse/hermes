'use strict';

function combineURLs(baseURL, relativeURL) {
    return relativeURL
        ? baseURL.replace(/\/+$/, '') + '/' + relativeURL.replace(/^\/+/, '')
        : baseURL;
}

function isFile(pathname) {
    return pathname.split('/').pop().indexOf('.') > -1;
}

function isDir(pathname) { return !isFile(pathname); }

exports.handler = (event, context, callback) => {

    // Extract the request from the CloudFront event that is sent to Lambda@Edge 
    var request = event.Records[0].cf.request;

    // Extract the URI from the request
    var olduri = request.uri;

    var newuri = olduri;

    if (!olduri.endsWith("index.html") && isDir(olduri)) {
        newuri = combineURLs(olduri, "index.html");
    }
/*
    // Log the URI as received by CloudFront and the new URI to be used to fetch from origin
    console.log("Old URI: " + olduri);
    console.log("New URI: " + newuri);
*/
    // Replace the received URI with the URI that includes the index page
    request.uri = newuri;

    // Return to CloudFront
    return callback(null, request);
};