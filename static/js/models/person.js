window.Person = Backbone.Model.extend({

    urlRoot: "people",

    initialize: function () {
        this.validators = {};

        this.validators.name = function (value) {
            return value.length > 0 ? {isValid: true} : {isValid: false, message: "You must enter a name"};
        };

        this.validators.assistant = function (value) {
            return value.match(/\s\.\s/gi) == null &&
                   value.match(/@/g) == null &&
                   value.match(/\s/g) == null
                    ? {isValid: true} : {isValid: false, message: "Please remove the domain extension and make sure this is a valid LDAP (first letter, last name)"};
        };

    },

    validateItem: function (key) {
        return (this.validators[key]) ? this.validators[key](this.get(key)) : {isValid: true};
    },

    // TODO: Implement Backbone's standard validate() method instead.
    validateAll: function () {

        var messages = {};

        for (var key in this.validators) {
            if(this.validators.hasOwnProperty(key)) {
                var check = this.validators[key](this.get(key));
                if (check.isValid === false) {
                    messages[key] = check.message;
                }
            }
        }

        return _.size(messages) > 0 ? {isValid: false, messages: messages} : {isValid: true};
    },

    defaults: {
        id: null,
        name: "",
        assistant: "",
        disabled: '0',
        is_manager: '0',
        mixers: {},
        assistant_for: {}
    }
});

window.PersonCollection = Backbone.Collection.extend({

    model: Person,
    initialize: function(props) { 
        this.groupname = props.groupname;
    },

    url: function() {
      return "group/" + this.groupname;
    }

});
