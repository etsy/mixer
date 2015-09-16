window.NewPersonView = Backbone.View.extend({

    initialize: function () {
        this.model.set('name', '');
        this.render();
    },

    render: function () {
        $(this.el).html(this.template(this.model.toJSON()));

        this.sharedView = new SharedPersonView({model: this.model});
        this.$('.shared_person_view').append(this.sharedView.el);

        return this;
    },

    events: {
        "change"        : "change",
        "click .save"   : "savePerson"
    },

    change: function (event) {
        // Remove any existing alert message
        utils.hideAlert();

        // Apply the change to the model
        var target = event.target;
        var change = {};
        change[target.name] = target.value;
        this.model.set(change);

        // Run validation rule (if any) on changed item
        var check = this.model.validateItem(target.id);
        if (check.isValid === false) {
            utils.addValidationError(target.id, check.message);
        } else {
            utils.removeValidationError(target.id);
        }
    },

    savePerson: function () {
        var self = this;

        var check = this.model.validateAll();
        if (check.isValid === false) {
            utils.displayValidationErrors(check.messages);
            return false;
        }

        var mixer_map = utils.getMixerMap(this.model);

        var person = new Person({name: this.model.get('name'),
                                 username: this.model.get('username'),
                                 assistant: this.model.get('assistant'),
                                 mixers: mixer_map
                                });

        person.save(null, {
            success: function (model) {
                app.navigate('people/' + model.id, true);
                utils.showAlert('Success!', 'Saved successfully', 'alert-success');
            },
            error: function () {
                utils.showAlert('Error', 'An error occurred while trying to save this item', 'alert-error');
            }
        });

        return false;
    }

});
