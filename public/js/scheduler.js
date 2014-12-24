$(function(){
  $("#addItem").click(function() {
    $("#schedulerItems").append(template);
  });

  $(".deleteItem").click(function(){
    $(this).parent().remove();
  });
});

var template = $("#row").html();
