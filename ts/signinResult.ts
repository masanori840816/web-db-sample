type SigninResult = {
	succeeded: boolean
	errorMessage: string
    nextUrl: string
}
// eslint-disable-next-line @typescript-eslint/no-explicit-any
export function isSigninResult(value: any): value is SigninResult {
    if(value == null) {
        return false;
    }
    return (("succeeded" in value) &&
        ("errorMessage" in value) &&
        ("nextUrl" in value) &&
        (typeof value.succeeded === "boolean") &&
        (typeof value.errorMessage === "string") &&
        (typeof value.nextUrl === "string"));
}