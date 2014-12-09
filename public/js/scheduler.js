$(function(){
  $("#addItem").click(function() {
    $("#schedulerItems").append(template);
  });

  $(".deleteItem").click(function(){
    $(this).parent().remove();
  })
});

var template = "    <div class=\"form-group row\">\
      <div class=\"col-sm-3\">\
        <input type=\"text\" class=\"form-control\" name=\"cronStrings[]\" value=\"\">\
      </div>\
      <div class=\"col-sm-3\">\
        <input type=\"text\" class=\"form-control\" name=\"funcNames[]\" value=\"\">\
      </div>\
    </div>"
