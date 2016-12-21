window.utils = {

    mixers : null,

    mixer_types : function() {
        this.mixers = [];
        var self = this;
        var p = $.ajax({
            url: 'mixers',
            dataType: "json",
            type: 'GET',
            processData: false,
            cache: true,
            async: true,
            contentType: false
        });
        p.done(function (mixer_data) {
            self.mixers = _.map(mixer_data, function(mixer) {
                return { name: mixer.name, description: mixer.description };
            });
        })
        p.fail(function () {
            self.showAlert('Error!', 'An error occurred while fetching mixers');
        });
        return p;
    },

    mixer_data : function(mixer) {
        var self = this;
        var m = $.ajax({
            url: 'mixer/' + mixer,
            dataType: "json",
            type: 'GET',
            processData: false,
            cache: true,
            async: true,
            contentType: false
        });
        m.fail(function() {
            self.showAlert('Error!', 'An error occurred while fetching data for ' + mixer);
        });
        return m;
    },

    isValidForUser: function(type, is_manager) {
        if ((type == "Managers" && is_manager == 1) || (type != "Managers")) {
            return true;
        }
        return false;
    },

    // Asynchronously load templates located in separate .html files
    loadTemplate: function(views, callback) {

        var deferreds = [];

        $.each(views, function(index, view) {
            if (window[view]) {
                deferreds.push($.get('tpl/' + view + '.html', function(data) {
                    window[view].prototype.template = _.template(data);
                }));
            } else {
                alert(view + " not found");
            }
        });

        $.when.apply(null, deferreds).done(callback);
    },

    displayValidationErrors: function (messages) {
        this.showAlert('Warning!', 'Fix validation errors and try again', 'alert-warning');
        for (var key in messages) {
            if (messages.hasOwnProperty(key)) {
                this.addValidationError(key, messages[key]);
            }
        }
    },

    spacesToDashes: function(word) {
        return word.replace(/\s/g, "-");
    },

    addValidationError: function (field, message) {
        var controlGroup = $('#' + field).parent().parent();
        $('.help-inline', controlGroup).html(message);
        $('.help-inline', controlGroup).addClass('alert alert-danger');
        $('.help-inline', controlGroup).prepend('<span class="glyphicon glyphicon-exclamation-sign" aria-hidden="true"></span>');
    },

    removeValidationError: function (field) {
        var controlGroup = $('#' + field).parent().parent();
        controlGroup.removeClass('error');
        $('.help-inline', controlGroup).html('');
    },

    showAlert: function(title, text, klass) {
        $('.alert').removeClass("alert-error alert-warning alert-success alert-info");
        $('.alert').addClass(klass);
        $('.alert').html('<strong>' + title + '</strong> ' + text);
        $('.alert').show();
    },

    // this is the map of mixers for a particular user, which is the model
    getMixerMap: function(model) {
        var mixer_map = model.get('mixers');
        if (mixer_map == null) {
            mixer_map = {};
        }
        $.each(this.mixer_types(), function( index, type ) {
            if (typeof mixer_map[type.name] == "undefined") {
                mixer_map[type.name] = model.get(type.name);
            }
        });
        return mixer_map;
    },

    hideAlert: function() {
        $('.alert').hide();
    }

};
