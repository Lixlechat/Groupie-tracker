function barrederecherche() {
    var VImput = document.getElementById('search').value;
    var VData = document.getElementById('searchs').options;
    console.log(VData)
    for(let i=0; i < VData.length; i++) {
        var VTransformeddata = Rtransformer(VData[i].value.toLowerCase());
        if(VData[i].value.toLowerCase()== VImput.toLowerCase() || VTransformeddata == VImput.toLowerCase()) {
            var ID = VData[i].dataset.id;
            window.location.href="/artist/"+ID;
            console.log(ID)
            return;
        }
    }

    // document.location.href="http://localhost:8010/error/";
    
}
console.log("fichierjscharger")

function Rtransformer(texte) {
    var Ntexte = "";
    for(i = 0; i < texte.length; i ++) {
        if(texte[i] == " " && texte[i+1] == "|") {
            return Ntexte;
        } else {
            Ntexte = Ntexte + texte[i];
        }
    }
}