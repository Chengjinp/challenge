function escapeHtml(str) {
    return str.replace(/&/g, "&").replace(/</g, "<").replace(/>/g, ">").replace(/"/g, "\"").replace(/'/g, "'");
}

function validateForm() {
    var tel = document.forms["phForm"]["tel"].value;
    if (tel == NaN || tel === '' || tel.match(/\d/g).length != 10) {
      alert("Telephone must be 10 digits only");
      return false;
    }
    return true;
  } 