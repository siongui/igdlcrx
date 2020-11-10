chrome.webRequest.onBeforeSendHeaders.addListener(
  function (details) {
    if (details.url.includes("://i.instagram.com/")) {
      for (var i = 0; i < details.requestHeaders.length; ++i) {
        if (details.requestHeaders[i].name === 'User-Agent') {
          details.requestHeaders[i].value = "Mozilla/5.0 (iPhone; CPU iPhone OS 13_3_1 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Mobile/15E148 Instagram 134.0.0.25.116 (iPhone12,5; iOS 13_3_1; en_US; en-US; scale=3.00; 1242x2688; 204771128)";
          break;
        }
      }
    }
    return { requestHeaders: details.requestHeaders };
  },
  { urls: ['<all_urls>'] },
  ['blocking', 'requestHeaders']
);
