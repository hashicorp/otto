(function(Sidebar, CubeDraw){

// Quick and dirty IE detection
var isIE = (function(){
	if (window.navigator.userAgent.match('Trident')) {
		return true;
	} else {
		return false;
	}
})();

// isIE = true;

var Init = {

	start: function(){
		var id = document.body.id.toLowerCase();

		if (this.Pages[id]) {
			this.Pages[id]();
		}

		//always init sidebar
		Init.initializeSidebar();
	},

	initializeSidebar: function(){
		new Sidebar();
	},

	initializeHomepage: function(){
		$('#use-case-nav a').click(function (e) {
		  e.preventDefault()
			console.log($(this)[0])
			console.log($(this).tab())
		  $(this).tab('show')
		})

		$('a[data-toggle="tab"]').on('show.bs.tab', function (e) {
		  console.log('show new active tab', e.target) // newly activated tab
		  console.log(e.relatedTarget);
		})
	},

	Pages: {
		'page-home': function(){
			Init.initializeHomepage();
		}
	}

};

Init.start();

})(window.Sidebar, window.CubeDraw);
