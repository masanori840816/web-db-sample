declare global {
    interface Window {
        SigninPage: SigninPageApi,
    };
}
export interface SigninPageApi {
    signin: (baseUrl: string) => void,
}