function startWiggleBar() {
        $('.wggl').hover(
                function() {
                        $(this).animate({ left: "-4px" }, 100);
                        $(this).animate({ left: "+4px" }, 100);
                        $(this).animate({ left: "+2px" }, 60);
                        $(this).animate({ left: "-2px" }, 60);
                        $(this).animate({ left: "+1px" }, 20);
                        $(this).animate({ left: "-1px" }, 20);
                        $(this).animate({ left: "0px" }, 10);
                },
                function() {
                        $(this).animate({ left: "0px" }, 100);
                }
        );
}

function typeset() {
}
function start() {
        startWiggleBar();
        //setTimeout('typeset()', 3000);
        //$.getScript("/s/mathjax/MathJax.js");
}

$(document).ready(start);
