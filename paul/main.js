//*compare la valeur a l'entrée avec les valeurs present dans la datalist pour nous renvoyer a la bonne page

function barrederecherche() {
    var VImput = document.getElementById('search').value;
    var VData = document.getElementById('searchs').options;
    console.log(VData)
    for (let i = 0; i < VData.length; i++) {
        var VTransformeddata = Rtransformer(VData[i].value.toLowerCase());
        if (VData[i].value.toLowerCase() == VImput.toLowerCase() || VTransformeddata == VImput.toLowerCase()) {
            var ID = VData[i].dataset.id;
            window.location.href = "/artist/" + ID;
            console.log(ID)
            return;
        }
    }

    window.location.href = "http://localhost:8010/error/";

}
console.log("fichierjscharger")

//* transforme les donnée present dans la datalist  pour les recherches taper a la manos. ça supprime le | pour la description
function Rtransformer(texte) {
    var Ntexte = "";
    for (i = 0; i < texte.length; i++) {
        if (texte[i] == " " && texte[i + 1] == "|") {
            return Ntexte;
        } else {
            Ntexte = Ntexte + texte[i];
        }
    }
}

function boutoncaché() {
    var otherCheckbox = document.querySelector('input[value="other"]');
    var otherText = document.querySelector('input[id="otherValue"]');
    otherText.style.visibility = 'hidden';

    otherCheckbox.onchange = function() {
        if (otherCheckbox.checked) {
            otherText.style.visibility = 'visible';
            otherText.value = '';
        } else {
            otherText.style.visibility = 'hidden';
        }
    };
}

/* When the user clicks on the button,
toggle between hiding and showing the dropdown content */
function myFunction() {
    document.getElementById("myDropdown").classList.toggle("show");
}

// Close the dropdown menu if the user clicks outside of it
window.onclick = function(event) {
    if (!event.target.matches('.dropbtn')) {
        var dropdowns = document.getElementsByClassName("dropdown-content");
        var i;
        for (i = 0; i < dropdowns.length; i++) {
            var openDropdown = dropdowns[i];
            if (openDropdown.classList.contains('show')) {
                openDropdown.classList.remove('show');
            }
        }
    }
}

function openNav() {
    document.getElementById("mySidebar").style.width = "300px";
    document.getElementById("main").style.marginLeft = "300px";
}

/* Set the width of the sidebar to 0 and the left margin of the page content to 0 */
function closeNav() {
    document.getElementById("mySidebar").style.width = "0";
    document.getElementById("main").style.marginLeft = "0";
}


