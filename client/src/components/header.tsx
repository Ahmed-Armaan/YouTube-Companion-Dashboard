import { useNavigate } from "react-router"
import Logo from "./logo"

function Header() {
	const navigate = useNavigate()

	const logOut = async () => {
		const backendURL = import.meta.env.VITE_BACKEND_URL
		const reqURL = new URL("/logout", backendURL)

		await fetch(reqURL.toString(), {
			method: "POST",
			credentials: "include",
		})

		localStorage.clear()
		sessionStorage.clear()

		navigate("/login")
	}

	return (
		<div className="flex justify-between p-3 bg-blue-200">
			<Logo subtitle={false} />
			<button onClick={logOut}>Logout</button>
		</div>
	)
}

export default Header
