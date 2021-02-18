//----------------------------DOMView--------------------------------------------

console.log( "DOMView.js Start");

/**
 * This wrapper makes it easier to handle the DOM View JS calls.
 * The actual calls for DOMView are:
 * view = new atv.DOMView()
 * view.onUnload - similar to onPageUnload
 * view.load ( XMLDOC, function(sucess) { ... } ) - pushes the view onto the stack the callback function is called back and gives you a success or fail call.
 * view.unload - removes the view from the stack.
 */
var DomViewManager = ( function() {
  var views = {},
      ViewNames = [],
      config = {},
      callbackEvents = {},
      optionDialogXML = '<?xml version="1.0" encoding="UTF-8"?> \
      <atv> \
        <body> \
          <optionDialog id="domview.optionDialog"> \
            <header> \
              <simpleHeader accessibilityLabel=""> \
                <title></title> \
              </simpleHeader> \
            </header> \
            <description></description> \
            <menu> \
              <initialSelection> \
                <row></row> \
              </initialSelection> \
              <sections> \
                <menuSection> \
                  <items> \
                  </items> \
                </menuSection> \
              </sections> \
            </menu> \
          </optionDialog> \
        </body> \
      </atv>';


  function _saveView( name, view ) {
    if( name && view ) {
      views[ name ] = view;
      _addViewToList( name );

    } else {
      console.error( "When saving a view, both name and view are required" );
    };
  };

  function _deleteView( name ) {
    if( views[ name ] ) {
      delete views[ name ];
      _removeViewFromList( name );
    };
  };

  function _retrieveView( name ) {
    if( name ) {
      return views[ name ] || null;
    } else {
      console.error( "When attempting to retrieve a view name is required.");
    };
    return null;
  };

  function _addViewToList( name ) {
    var index = ViewNames.indexOf( name );
    if( index == -1 ) {
      ViewNames.push( name );
    };
  };

  function _removeViewFromList( name ) {
    var index = ViewNames.indexOf( name );
    if( index > -1 ) {
      ViewNames.splice( index, 1 );
    };
  };

  function _createDialogXML( dialogOptions ) {
    var doc = atv.parseXML( optionDialogXML ),
        title = dialogOptions.title,
        description = dialogOptions.description,
        initialSelection = dialogOptions.initialSelection || 0,
        options = dialogOptions.options || [];


    // fill in the title, accessibility label
    doc.rootElement.getElementByTagName( 'title' ).textContent = title;
    doc.rootElement.getElementByTagName( 'simpleHeader' ).setAttribute( 'accessibilityLabel', title +". "+ description );

    // fill in the description
    doc.rootElement.getElementByTagName( 'description' ).textContent = description;

    // fill in the initial selection
    doc.rootElement.getElementByTagName( 'row' ).textContent = initialSelection;

    // fill in the options
    var items = doc.rootElement.getElementByTagName( 'items' );
    options.forEach( function ( option, index ) {
      // save option callbacks
      RegisterCallbackEvent( "DialogOption_"+index, option.callback );

      // create the option
      var newOptionButton = ATVUtils.createNode({
          "name": "oneLineMenuItem",
          "attrs": [{
            "name": "id",
            "value": "DialogOption_"+ index
          }, {
            "name": "accessibilityLabel",
            "value": option.label
          }, {
            "name": "onSelect",
            "value": "DomViewManager.fireCallback( 'DialogOption_"+ index +"' );"
          }],
          "children": [{
            "name": "label",
            "text": option.label
          }]
        },
        doc );

      // append it to the items.
      items.appendChild( newOptionButton );
    });

    return doc;

  }

  function ListSavedViews() {
    return ViewNames;
  };

  function setConfig(property, value) {
    console.log( " ===> Setting: "+ property +" = "+ value +" <=== " );
    config[ property ] = value;
  };

  function getConfig(property) {
    var value = config[property];
    return (value) ? value: null;
  };

  // Create a new DomView
  function CreateView( name, dialogOptions ) {
    if( name ) {
      var view = new atv.DOMView();

      _saveView( name, view );

      if( typeof( dialogOptions ) === "object" ) {
        var doc = _createDialogXML( dialogOptions );
      };

      setConfig( name+"_doc", doc )

      view.onUnload = function() {
        console.log(" == DOMView onUnload called == " );
        FireCallbackEvent("ONUNLOADVIEW", {
          "name": name,
          "view": this
        });
      };

    } else {
      console.error("When attempting to create a DOM view, name is required.");
    };
  };

  function RemoveView( name ) {
    // unload the view, remove the view from the view list, remove the view name
    UnloadView( name );
    _deleteView( name );
  };

  function LoadView( name, doc ) {
    try {
      var view = _retrieveView( name ),
          doc = doc || getConfig( name+"_doc" );

      if( !view )
      {
        CreateView( name );
        view = _retrieveView( name );
      }

      console.log( "We load the view: "+ name +" : "+ view );
      view.load(doc, function(success) {
          console.log("DOMView succeeded " + success);
          if( success )
          {
            console.log("=== Saving Document: "+ name +"_doc ===");
            view.doc = doc.serializeToString();
            FireCallbackEvent( "ONLOADSUCCESS", { "view": name } )
          }
          else
          {
            var msg = "Unable to load view."
            FireCallbackEvent( "ONLOADERROR", { "id": "LOADERROR", "view":name, "msg": msg } );
          }
      });
    } catch ( error ) {
      console.error( "LOAD ERROR: "+ error );
    };
  };

  function UnloadView( name ) {
    var view = _retrieveView( name );
    view.unload();
  };

  function RegisterCallbackEvent( name, callback ) {
    console.log(" ---- Registering Callback: " + name + " with callback type: " + typeof(callback));
    if (typeof callback === "function") {
      callbackEvents[name] = callback;
    } else {
      console.error("When attempting to register a callback event, a callback function is required.");
    };
  };

  function FireCallbackEvent( name, parameters, scope ) {
    var scope = scope || this,
    parameters = parameters || {};

    if (callbackEvents[name] && typeof callbackEvents[name] === "function") {
      callbackEvents[name].call(scope, parameters)
    };
  };

  return {
    "createView": CreateView,
    "removeView": RemoveView,
    "loadView": LoadView,
    "unloadView": UnloadView,
    "listViews": ListSavedViews,
    "registerCallback": RegisterCallbackEvent,
    "fireCallback": FireCallbackEvent
  };

})();


// ------ End DOM View Manager --------

console.log("DomView.js end");




// ------ Usage function --------

function optionOne() {
    console.log("Option one taken. " + this);
    console.log( JSON.stringify( DomViewManager.listViews() ) );
    DomViewManager.unloadView( "DialogView" );
}

function optionTwo() {
    console.log("Option two taken");
    DomViewManager.unloadView( "DialogView" );
}


var optionDialogXML = '<?xml version="1.0" encoding="UTF-8"?> \
<atv> \
  <body> \
    <optionDialog id="com.sample.error-dialog"> \
      <header> \
        <simpleHeader accessibilityLabel="Dialog with Options"> \
          <title>Option Dialog</title> \
        </simpleHeader> \
      </header> \
      <description>Lorem ipsum dolor sit amet, consectetur adipisicing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.</description> \
      <menu> \
        <initialSelection> \
          <row>1</row> \
        </initialSelection> \
        <sections> \
          <menuSection> \
            <items> \
              <oneLineMenuItem id="option1" accessibilityLabel="Option 1" onSelect="optionOne()"> \
                <label>Option 1</label> \
              </oneLineMenuItem> \
              <oneLineMenuItem id="option2" accessibilityLabel="Option 2" onSelect="optionTwo()"> \
                <label>Option 2</label> \
              </oneLineMenuItem> \
            </items> \
          </menuSection> \
        </sections> \
      </menu> \
    </optionDialog> \
  </body> \
</atv>';


function LoadDialogDOMView() {
  var doc = atv.parseXML( optionDialogXML );

  console.log( "Creating the view --> " );
  DomViewManager.createView( "DialogView" );

  console.log( "Loading the view --> " );
  DomViewManager.loadView( "DialogView", doc );

  console.log( "View is loaded" );
}

//-------------------------------------------------------------------------------

console.log("DomView.js end");
