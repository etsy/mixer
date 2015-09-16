window.PersonView = Backbone.View.extend({

    initialize: function () {
        this.render();
    },

    render: function () {
        $(this.el).html(this.template(this.model.toJSON()));

        _.bindAll(this, "savePerson");
        this.sharedView = new SharedPersonView({model: this.model, saveCallback: this.savePerson});
        this.$('.shared_person_view').append(this.sharedView.el);

        return this;
    },

    events: {
        "change"        : "change",
    },

    change: function (event) {
        // Remove any existing alert message
        utils.hideAlert();

        // Apply the change to the model
        var target = event.target;
        if (target.type != "checkbox") {

            var change = {};
            change[target.name] = target.value;
            this.model.set(change);

            // Run validation rule (if any) on changed item
            var check = this.model.validateItem(target.id);
            if (check.isValid === false) {
                utils.addValidationError(target.id, check.message);
                return false;
            } else {
                utils.removeValidationError(target.id);
            }

            this.model.save();
        }

    },

    savePerson: function () {
        var self = this;
        var check = this.model.validateAll();

        var mixer_map = utils.getMixerMap(this.model);
        this.model.set('mixers', mixer_map);

        if (check.isValid === false) {
            utils.displayValidationErrors(check.messages);
            return false;
        }

        this.model.save(null, {
            success: function (model) {
                utils.showAlert('Success!', 'Saved successfully', 'alert-success');
            },
            error: function () {
                utils.showAlert('Error', 'An error occurred while trying to save', 'alert-error');
            }
        });
    },

});
