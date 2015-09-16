var AppRouter = Backbone.Router.extend({

    routes: {
        ""                       : "home",
        "group/:name"            : "list",
        "group/:name/page/:page" : "list",
        "people/:id"             : "peopleDetails",
        "about"                  : "about"
    },

    initialize: function () {
        utils.mixer_types().done(function(types) {
            _.each(utils.mixers, function(type){
                $("#nav-header").append('<li class="'+ type +'-menu"><a href="#group/'+ type +'">'+ type +'</a></li>');
            });
        });
        this.headerView = new HeaderView();
        $('.header').html(this.headerView.el);
    },

    home: function () {
        var auth = new AuthUser();

        auth.fetch({
            success: function () {
                if (auth.id == 0) {
                    // display message here to sign up
                    $("#content").html(new NewPersonView({model: auth}).el);
                } else {
                    app.navigate('people/' + auth.id, true);
                }
            },
            error: function () {
                console.log('couldnt auth...');
            }
        });
    },

	list: function(gn, page) {
        var p = page ? parseInt(page, 10) : 1;
        var personList = new PersonCollection({groupname: gn});
        personList.fetch({success: function(){
            $("#content").html(new PersonListView({model: personList, page: p}).el);
        }});
        this.headerView.selectMenuItem(gn + '-menu');
    },

    authenticate: function(person) {
        var auth = new AuthUser();

        auth.fetch({
            success: function () {
                person.set("disabled", '0');
                if ((auth.id != person.id) && (auth.get('username') != person.get('assistant'))) {
                    person.set("disabled", '1');
                }
            },
            error: function () {
                console.log('couldnt auth...');
            }
        });
    },

    peopleDetails: function (id) {
        var person = new Person({id: id});

        var self = this;
        person.fetch({success: function(){
            self.authenticate(person);
            $("#content").html(new PersonView({model: person}).el);
        }});
        this.headerView.selectMenuItem('me-menu');
    },

    about: function () {
        if (!this.aboutView) {
            this.aboutView = new AboutView();
        }
        $('#content').html(this.aboutView.el);
        this.headerView.selectMenuItem('about-menu');
    }

});

utils.loadTemplate(
[
    'HeaderView',
    'PersonView',
    'SharedPersonView',
    'NewPersonView',
    'PersonListItemView',
    'AboutView'
],
function() {
    app = new AppRouter();
    Backbone.history.start();
});
