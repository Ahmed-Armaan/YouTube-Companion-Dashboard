import { useEffect } from 'react'
import { useNavigate } from 'react-router'
import GoogleButton from 'react-google-button'
import Logo from '../components/logo'
import '../index.css'

function Login() {
  const navigate = useNavigate()

  const loginWithGoogle = () => {
    var url = new URL("https://accounts.google.com/o/oauth2/v2/auth")
    const client_id = import.meta.env.VITE_GOOGLE_CLIENT_ID
    const redirect_url = import.meta.env.VITE_REDIRECT_URI
    const scopes = import.meta.env.VITE_GOOGLE_OAUTH_SCOPE

    url.searchParams.append("client_id", client_id)
    url.searchParams.append("redirect_uri", redirect_url)
    url.searchParams.append("response_type", "code")
    url.searchParams.append("scope", scopes)
    url.searchParams.append("access_type", "offline")
    url.searchParams.append("prompt", "consent")

    window.location.href = url.toString()
  }

  const checkLoggedIn = async () => {
    const backendURL = import.meta.env.VITE_BACKEND_URL
    const reqPath = "/me"
    const reqURL = new URL(backendURL, reqPath)

    const res = await fetch(reqURL, {
      credentials: "include",
    })

    if (res.status === 200) {
      navigate("/home")
    }
  }

  useEffect(() => {
    checkLoggedIn()
  }, [])

  return (
    <>
      {/* <button className="google-btn" onClick={loginWithGoogle}>
          Login with Google
        </button> */}

      <div className="flex h-screen w-full items-center justify-center">
        <div className="flex flex-col items-center gap-4 p-5 border-black border-2 rounded-xl">
          <div className="text-xl text-center">
            <Logo subtitle />
          </div>

          <GoogleButton onClick={loginWithGoogle} />
        </div>
      </div>
    </>
  )
}

export default Login 
