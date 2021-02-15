console.log( "Loading the testing code" );

function changeMyLabel( itemid ) {
	var ele = document.getElementById( itemid ),
		newLabel = "Updated: "+ Date.now();

	console.log( "the new label for : "+ ele.tagName +" is: "+ newLabel );

	ele.getElementByTagName( 'label' ).textContent = newLabel;
}

function changeTheTitle( ) {
	var title = document.evaluateXPath( '//header//title', document ),
		newTitle = Date.now();

	console.log( "the new label for : "+ title[0].tagName +" with current title : "+ title[0].textContent +" is: "+ newTitle );

	title[0].textContent = newTitle;

	console.log( "the label for : "+ title[0].tagName +" is now "+ title[0].textContent );
}

function changeMyTitle( itemid ) {
	var ele = document.getElementById( itemid ),
		newLabel = "Updated: "+ Date.now();

	console.log( "the new label for : "+ ele.tagName +" is: "+ newLabel );

	ele.getElementByTagName( 'title' ).textContent = newLabel;
}

function handleNavbarNavigate( event ) {
	console.log( "Handling the navigation event."+ JSON.stringify( event ) );

		// The navigation item ID is passed in through the event parameter.
	var navId = event.navigationItemId,

		// Use the event.navigationItemId to retrieve the appropriate URL information this can
		// retrieved from the document navigation item.
		docUrl = document.getElementById( navId ).getElementByTagName( 'url' ).textContent,

		// Request the XML document via URL and send any headers you need to here.
		ajax = new ATVUtils.Ajax({
			"url": docUrl,
			"success": function( xhr ){
				console.log( "successfully loaded the XHR" );

				// After successfully retrieving the document you can manipulate the document
				// before sending it to the navigation bar.
				var doc = xhr.responseXML,
					title = doc.rootElement.getElementByTagName( 'title' );
				// title.textContent = title.textContent +": Appended by Javascript";

				// Once the document is ready to load pass it to the event.success function
				event.success( doc );
			},
			"failure": function( status, xhr ){
				// If the document fails to load pass an error message to the event.failure button
				event.failure( "Navigation failed to load." );
			}
		});

	event.onCancel = function() {
		console.log("nav bar nagivation was cancelled");
		// declare an onCancel handler to handle cleanup if the user presses the menu button before the page loads.
	}

}

function atvStressTest() {
	console.log( "Loading the stress test" );
	ATVUtils.loadURL({
		"url": "http://sample-web-server/sample-xml/k66-stress-test.xml",
		"processXML": function ( doc ) {
			console.log( "Stress test is loaded: --> " );
			var movies = doc.rootElement.getElementsByTagName( 'moviePoster' )
				i = 0;

			console.log( "Array of movie posters has been created: length: "+ movies.length +" : commence looping --> " );

			while( movies.length > 250 ) {
				var movie = movies.pop();

				movie.removeFromParent();

			};

			console.log( " <-- I have looped through the movie posters and limited the size to 250 " );
		}
	});
}

function localTimeTest() {
	console.log( "In this function we will be printing out various versions of local time.");
	var now = new Date(),
		standardFormats = {
			"era": 'GGG GGGG GGGGG',
			"year": 'y yy yyy yyyy Y YY YYY YYYY u uu uuu uuuu U UU UUU UUUU',
			"quarter": 'qq qqq qqqq QQ QQQ QQQQ',
			"month": 'MM MMM MMMM MMMMM LL LLL LLLL LLLLL',
			"week": 'w ww W',
			"day": 'd D F gggggg',
			"weekDay": 'EEE EEEE EEEEE ee eee eeee eeeee c ccc cccc ccccc',
			"period": 'a',
			"hour": 'hh HH kk KK',
			"minute": 'mm',
			"second": 'ss SSSSS AAAAAAAAAA',
			"zone": 'zzz zzzz ZZZ ZZZZ ZZZZZ v vvvv V VVVV'

		},
		examples = [
			"MMddYYYY",
			"hhmmss",
			"HHmmss",
			"EEE MMMM dd YYYY",
			"EEEMMMMddYYYY"
		];

	console.log( " == Local Time Examples == " );
	for( var p in standardFormats )
	{
		if( standardFormats.hasOwnProperty( p ) )
		{
			var format = standardFormats[ p ];
			console.log( " - "+ p +" formats for pattern '"+ format +"' - ");
			format.split( ' ' ).forEach( function( f ) {
				console.log( " ----> '"+ f +"': "+ atv.localtime( now, f ) +" <----- " );
			});
			console.log( "\n" );
		}
	}

	console.log( " - EXAMPLES - " );
	examples.forEach( function( example ) {
		console.log( " ----> '"+ example +"': "+ atv.localtime( now, example ) +" <----- " );
	});
	console.log( " == End Local Time Examples == \n\n" );

}

localTimeTest();