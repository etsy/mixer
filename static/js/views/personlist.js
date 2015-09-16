window.PersonListView = Backbone.View.extend({

    initialize: function () {
        this.render();
    },

    render: function () {
        var people = this.model.models;
        var len = people.length;
        var startPos = (this.options.page - 1) * 12;
        var endPos = Math.min(startPos + 12, len);

        $(this.el).html('<ul class="thumbnails col-md-12"></ul>');

        for (var i = startPos; i < endPos; i++) {
            $('.thumbnails', this.el).append(new PersonListItemView({model: people[i]}).render().el);
        }

        $(this.el).append(new Paginator({model: this.model, page: this.options.page}).render().el);

        return this;
    }
});

window.PersonListItemView = Backbone.View.extend({

    tagName: "li",

    className: "col-md-3",

    initialize: function () {
        this.model.bind("change", this.render, this);
        this.model.bind("destroy", this.close, this);
    },

    render: function () {
        $(this.el).html(this.template(this.model.toJSON()));
        return this;
    }

});
