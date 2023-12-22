export const  copyToClipboard = (text) => {
    if(navigator.clipboard) {
        navigator.clipboard.writeText(text);
    }
    else{
        alert(text);
    }
}