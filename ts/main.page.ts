import { isSigninResult } from "./signinResult";

window.MainPage = {
    signout(baseUrl: string) {
        fetch(`${baseUrl}signout`, {
            method: "GET",
        })
        .then(res => res.json())
        .then(j => {
            if(isSigninResult(j)) {
                if(j.succeeded === false) {
                    alert(j.errorMessage);
                    return;
                }
                location.href = j.nextUrl;
            } else {
                alert("Sign out error");
            }
        })
        .catch(err => console.error(err));        
    }
}