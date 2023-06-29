declare global {
    interface Window {
        SigninPage: SigninPageApi,
        MainPage: MainPageApi,
    };
}
export interface SigninPageApi {
    signin: (baseUrl: string) => void,
}
export interface MainPageApi {
    signout: (url: string) => void,
}