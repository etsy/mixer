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
                if (utils.isValidForUser(type.name, self.model.get('is_manager'))) {

                    var mixerContainer = $("<div/>", {
                        class: "mixer-container",
                    });

                    var infoWrapper = $("<div/>", {
                        class: "mixer-information",
                    });

                    var infoLabel = $("<span/>", {
                        class: "mixer-label",
                        text: type.name
                    });

                    var infoDescription = $("<span/>", {
                        class: "mixer-description",
                        text: type.description
                    });

                    var checkboxWrapper = $("<div/>", {
                        class: "mixer-checkbox",
                    });

                    var checkboxInput = $("<input/>", {
                        type: "checkbox",
                        class: "mixers",
                        name: type.name,
                        value: "1",
                        disabled: (self.model.get('disabled') === '1' ? true : false)
                    });
                    checkboxInput.prop('checked', self.model.get('mixers')[type.name]);

                    infoWrapper
                    .append(infoLabel)
                    .append(infoDescription);

                    checkboxWrapper.append(checkboxInput);

                    mixerContainer
                    .append(infoWrapper)
                    .append(checkboxWrapper)
                    .appendTo("#participating-mixers");

                    checkboxInput.bootstrapSwitch('onText', 'YES');
                    checkboxInput.bootstrapSwitch('offText', 'NO');
                    checkboxInput.on("switchChange.bootstrapSwitch", $.proxy(self.handleMixerChange, self));
                }
            });
        });

        return this;
    }

});
