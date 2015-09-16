window.SharedPersonView = Backbone.View.extend({

    tagName: "div",

    handleMixerChange : function (e) {
        var mixers = this.model.get('mixers');
        mixers[e.target.name] = e.target.checked;
        this.model.set('mixers', mixers);

        if (this.options.saveCallback) {
            this.options.saveCallback();
        }
    },

    initialize: function () {
        this.model.bind("change", this.render, this);
        this.model.bind("destroy", this.close, this);
        this.render();
    },

    render: function () {
        $(this.el).html(this.template(this.model.toJSON()));
        var self = this;

        utils.mixer_types().done(function(types) {
            _.each(utils.mixers, function(type){
                if (utils.isValidForUser(type, self.model.get('is_manager'))) {

                    $("<span/>", {
                        class: "mixer-label",
                        text: type
                    })
                    .appendTo("#participating-mixers");

                    var input = $("<input/>", {
                        type: "checkbox",
                        class: "mixers",
                        name: type,
                        value: "1",
                        disabled: (self.model.get('disabled') === '1' ? true : false)
                    });
                    input.prop('checked', self.model.get('mixers')[type]);

                    $("<div/>", {
                        class: "mixer-container",
                    })
                    .append(input)
                    .append($("<br/>"))
                    .appendTo("#participating-mixers");

                    input.bootstrapSwitch('onText', 'YES');
                    input.bootstrapSwitch('offText', 'NO');
                    input.on("switchChange.bootstrapSwitch", $.proxy(self.handleMixerChange, self));
                }
            });
        });

        return this;
    }

});
