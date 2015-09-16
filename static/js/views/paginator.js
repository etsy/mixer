window.Paginator = Backbone.View.extend({

    className: "text-center",

    initialize:function () {
        this.model.bind("reset", this.render, this);
        this.render();
    },

    render:function () {

        var items = this.model.models;
        var len = items.length;
        var pageCount = Math.ceil(len / 12);

        $(this.el).html('<ul />');
        $('ul', this.el).addClass('pagination');

        for (var i=0; i < pageCount; i++) {
            $('ul', this.el).append("<li" + ((i + 1) === this.options.page ? " class='active'" : "") + "><a href='#group/"+ this.model.groupname +"/page/"+(i+1)+"'>" + (i+1) + "</a></li>");
        }

        return this;
    }
});
