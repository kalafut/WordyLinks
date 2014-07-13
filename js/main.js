var activeButton = "btn-sel";
var inactiveButton = "btn-nonsel";
var errorClass = "has-error";


$(".mode-btn").click(function() {
    if( !$(this).hasClass(activeButton) ) {
        $(".mode-btn").toggleClass(activeButton + " " + inactiveButton);
        $(".form-group").removeClass(errorClass);
        $(".help-inline").text("");

        setupForm(200);
    }
    return false;
});

function setupForm(duration) {
    if($("#save-btn").hasClass(activeButton)) {
        if ($('#link-ctrl').css('display') === 'none') {
            $('#link-ctrl').slideDown(duration);
        }
        $('#submit_btn').text("Save");
        $('#submit_form').attr('action', 'save');
        $('#url').focus();

        setTimeout(function() { $("#success_msg").css("visibility", "hidden"); }, 2000);
    } else {
        $('#link-ctrl').slideUp(duration);
        $('#submit_btn').text("Load");
        $('#submit_form').attr('action', 'load');
        $('#keywords').focus();
    }
}

setupForm(0);

$("#cover").hide();

$('#keywords, #url').keyup(function (e) {
    if (e.keyCode == 13) {
        $('#submit_form').submit();
    }
});

