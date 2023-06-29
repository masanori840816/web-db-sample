import { isSigninResult } from "./signinResult";

window.SigninPage = {
    signin(baseUrl: string) {
        const userName = document.getElementById("signin_user_name") as HTMLInputElement;
        const password = document.getElementById("signin_password_input") as HTMLInputElement;

        fetch(`${baseUrl}signin`,{
            method: "POST",
            body: JSON.stringify({
                userName: userName.value, 
                password: password.value
            }),
        })
        .then(res => res.json())
        .then(j => {
            const errorMessage = document.getElementById("signin_error_message") as HTMLElement;
            if(isSigninResult(j)) {
                if(j.succeeded === false) {
                    errorMessage.textContent = j.errorMessage;
                    return;
                }
                location.href = j.nextUrl;
            } else {
                errorMessage.textContent = "Sign in error";
            }
        })
        .catch(err => console.error(err));
    }
}