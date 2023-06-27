window.SigninPage = {
    signin(baseUrl: string) {
        const userName = document.getElementById("signin_user_name") as HTMLInputElement;
        const password = document.getElementById("signin_password_input") as HTMLInputElement;

        fetch(`${baseUrl}/signin`,{
            method: "POST",
            body: JSON.stringify({
                userName: userName.value, 
                password: password.value
            }),
        })
        .then(res => res.json())
        .then(json => console.log(json))
        .catch(err => console.error(err));
    }
}