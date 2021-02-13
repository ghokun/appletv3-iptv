/**
 * This is an XHR handler. It handles most of tediousness of the XHR request
 * and keeps track of onRefresh XHR calls so that we don't end up with multiple
 * page refresh calls.
 *
 * You can see how I call it on the handleRefresh function below.
 *
 *
 * @params object (hash) $options
 * @params string $options.url - url to be loaded
 * @params string $options.method - "GET", "POST", "PUT", "DELTE"
 * @params bool $options.type - false = "Sync" or true = "Async" (You should always use true)
 * @params func $options.success - Gets called on readyState 4 & status 200
 * @params func $options.failure - Gets called on readyState 4 & status != 200
 * @params func $options.callback - Gets called after the success and failure on readyState 4
 * @params string $options.data - data to be sent to the server
 * @params bool $options.refresh - Is this a call from the onRefresh event.
 */
ATVUtils.Ajax = function($options) {
	var me = this;
	$options = $options || {}

	/* Setup properties */
	this.url = $options.url || false;
	this.method = $options.method || "GET";
	this.type = ($options.type === false) ? false : true;
	this.success = $options.success || null;
	this.failure = $options.failure || null;
	this.data = $options.data || null;
	this.complete = $options.complete || null;
	this.refresh = $options.refresh || false;

	if(!this.url) {
		console.error('\nAjax Object requires a url to be passed in: e.g. { "url": "some string" }\n')
		return undefined;
	};

	this.id = Date.now();

	this.createRequest();

	this.req.onreadystatechange = this.stateChange;

	this.req.object = this;

	this.open();

	this.send();

};

ATVUtils.Ajax.currentlyRefreshing = false;
ATVUtils.Ajax.activeRequests = {};

ATVUtils.Ajax.prototype = {
	stateChange: function() {
		var me = this.object;
		switch(this.readyState) {
			case 1:
				if(typeof(me.connection) === "function") me.connection(this, me);
				break;
			case 2:
				if(typeof(me.received) === "function") me.received(this, me);
				break;
			case 3:
				if(typeof(me.processing) === "function") me.processing(this, me);
				break;
			case 4:
				if(this.status == "200") {
					if(typeof(me.success) === "function") me.success(this, me);
				} else {
					if(typeof(me.failure) === "function") me.failure(this.status, this, me);
				}
				if(typeof(me.complete) === "function") me.complete(this, me);
				if(me.refresh) Ajax.currentlyRefreshing = false;
				break;
			default:
				console.log("I don't think I should be here.");
				break;
		}
	},
	cancelRequest: function() {
		this.req.abort();
		delete ATVUtils.Ajax.activeRequests[ this.id ];
	},
	cancelAllActiveRequests: function() {
		for ( var p in ATVUtils.Ajax.activeRequests ) {
			if( ATVUtils.Ajax.activeRequests.hasOwnProperty( p ) ) {
				var obj = ATVUtils.Ajax.activeRequests[ p ];
				if( ATVUtils.Ajax.prototype.isPrototypeOf( obj ) ) {
					obj.req.abort();
				};
			};
		};
		ATVUtils.Ajax.activeRequests = {};
	},
	createRequest: function() {
		try {
			this.req = new XMLHttpRequest();
			ATVUtils.Ajax.activeRequests[ this.id ] = this;
			if(this.refresh) ATVUtils.Ajax.currentlyRefreshing = true;
		} catch (error) {
			alert("The request could not be created: </br>" + error);
			console.error("failed to create request: " +error);
		}
	},
	open: function() {
		try {
			this.req.open(this.method, this.url, this.type);
		} catch(error) {
			console.log("failed to open request: " + error);
		}
	},
	send: function() {
		var data = this.data || null;
		try {
			this.req.send(data);
		} catch(error) {
			console.log("failed to send request: " + error);
		}
	}
};