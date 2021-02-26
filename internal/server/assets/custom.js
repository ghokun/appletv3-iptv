var navbarItemNumber = null;

function updatePage(url) {
  if (navbarItemNumber == '0') // First navbar item is a special case
  {
    atv.loadAndSwapURL(url);
  } else {
    var req = new XMLHttpRequest();
    req.onreadystatechange = function () {
      if (req.readyState == 4) {
        var doc = req.responseXML;
        var navBar = doc.getElementById('templates/main.xml');
        var navKey = navBar.getElementByTagName('navigation');
        navKey.setAttribute('currentIndex', navbarItemNumber);
        atv.loadAndSwapXML(doc);
      }
    };
  }
  req.open('GET', url, false);
  req.send();
};

// Main navigation bar handler
function handleNavbarNavigate(event) {
  // The navigation item ID is passed in through the event parameter.
  var navId = event.navigationItemId;
  var root = document.rootElement;
  var navitems = root.getElementsByTagName('navigationItem')
  for (var i = 0; i < navitems.length; i++) {
    if (navitems[i].getAttribute('id') == navId) {
      navbarItemNumber = i.toString();
      break;
    }
  }
  // Use the event.navigationItemId to retrieve the appropriate URL information this can
  // retrieved from the document navigation item.
  docUrl = document.getElementById(navId).getElementByTagName('url').textContent,
    // Request the XML document via URL and send any headers you need to here.
    ajax = new ATVUtils.Ajax({
      "url": docUrl,
      "success": function (xhr) {
        // After successfully retrieving the document you can manipulate the document
        // before sending it to the navigation bar.
        var doc = xhr.responseXML,
          title = doc.rootElement.getElementByTagName('title');
        // title.textContent = title.textContent +": Appended by Javascript";
        // Once the document is ready to load pass it to the event.success function
        event.success(doc);
      },
      "failure": function (status, xhr) {
        // If the document fails to load pass an error message to the event.failure button
        event.failure("Navigation failed to load.");
      }
    });
  event.onCancel = function () {
    // declare an onCancel handler to handle cleanup if the user presses the menu button before the page loads.
  }
}

function callUrlAndUnload(url, method) {
  ajax = new ATVUtils.Ajax({
    "url": url,
    "method": method,
    "success": function (xhr) {
      atv.unloadPage();
    },
    "failure": function (status, xhr) {
      atv.unloadPage();
    }
  });
}

function addSpinner(elem) {
  var elem_add = document.makeElementNamed("spinner");
  elem.getElementByTagName("accessories").appendChild(elem_add);
}

function removeSpinner(elem) {
  var elem_remove = elem.getElementByTagName("accessories").getElementByTagName("spinner");
  if (elem_remove) elem_remove.removeFromParent();
}

function callUrlAndUpdateElement(id, url, method) {
  element = document.getElementById(id);
  rightLabel = element.getElementByTagName("rightLabel");
  if (rightLabel.textContent == '0') {
    return;
  }
  addSpinner(element);
  ajax = new ATVUtils.Ajax({
    "url": url,
    "method": method,
    "success": function (xhr) {
      removeSpinner(element);
      var doc = xhr.responseXML;
      newElement = doc.getElementById(id);
      rightLabel.textContent = newElement.getElementByTagName("rightLabel").textContent;
      if (rightLabel.textContent == '0') {
        element.setAttribute("dimmed", "true");
      }
    },
    "failure": function (status, xhr) {
      removeSpinner(element);
      rightLabel.textContent = '⚠️';
    }
  });
}

function editM3UAddress(title, instructions, label, footnote, defaultValue) {
  var textEntry = new atv.TextEntry();
  textEntry.type = 'emailAddress';
  textEntry.title = title;
  textEntry.instructions = instructions;
  textEntry.label = label;
  textEntry.footnote = footnote;
  textEntry.defaultValue = defaultValue;
  textEntry.defaultToAppleID = false;
  // textEntry.image = 'http://sample-web-server/sample-xml/images/ZYXLogo.png'
  var label2 = document.getElementById("edit-m3u").getElementByTagName('label2');
  textEntry.onSubmit = function (value) {
    ajax = new ATVUtils.Ajax({
      "url": "https://appletv.redbull.tv/set-m3u.xml?m3u=" + value,
      "method": "POST",
      "success": function (xhr) {
        label2.textContent = value;
      },
      "failure": function (status, xhr) {
        label2.textContent = status;
      }
    });
  }
  textEntry.onCancel = function () {
    label2.textContent = defaultValue;
  }
  textEntry.show();
}