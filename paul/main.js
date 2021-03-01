function barrederecherche() {
    var VImput = document.getElementById('search').value;
    var VData = document.getElementById('searchs').options;
    console.log(VData)
    for(let i=0; i < VData.length; i++) {
        var VTransformeddata = Rtransformer(VData[i].value.ToLowerCase());
        if(VData[i].value.ToLowerCase()== VImput.ToLowerCase() || VTransformeddata == VImput.ToLowerCase()) {
            var ID = VData[i].dataset.id;
            document.location.replace("http://localhost:8010/artist/"+ID);
            console.log(ID)
            return;
        }
    }

}
console.log("fichierjscharger")

function Rtransformer(texte) {
    var Ntexte = "";
    for(i = 0; i < text.length; i ++) {
        if(texte[i] == " " && texte[i+1] == "|") {
            return Ntexte;
        } else {
            Ntexte = Ntexte + texte[i];
        }
    }
}