import React, { useState, useEffect } from "react"
import { useNavigate } from "react-router-dom";
import { Endpoints } from "../api"
import { deleteCookie } from "../utils"
import Errors from "../components/Errors"

const Session = () => {
  const [user, setUser] = useState(null)
  const [isFetching, setIsFetching] = useState(false)
  const [errors, setErrors] = useState([])
  let navigate = useNavigate();

  const headers = {
    Accept: "application/json",
    Authorization: document.cookie.split("token=")[1],
  }

    const [data, setData] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        const getData = async () => {
        try {
            const response = await fetch(
            `http://localhost:9808`
            );
            if (!response.ok) {
            throw new Error(
                `This is an HTTP error: The status is ${response.status}`
            );
            }
            let actualData = await response.json();
            setData(actualData);
            setError(null);
        } catch(err) {
            setError(err.message);
            setData(null);
        } finally {
            setLoading(false);
        }  
        }
        getData()
    }, [])

  const getUserInfo = async () => {
    try {
      setIsFetching(true)
      const res = await fetch(Endpoints.session, {
        method: "GET",
        //credentials: "include",
        headers,
      })
    
      if (!res.ok) logout()

      const { success, errors = [], user } = await res.json()
      setErrors(errors)
      
      //if (!success) navigate("/v1/user/register")
      //setUser(user)
    } catch (e) {
      setErrors([e.toString()])
    } finally {
      setIsFetching(false)
    }
  }

  const logout = async () => {
    const res = await fetch(Endpoints.logout, {
      method: "GET",
      credentials: "include",
      headers,
    })

    if (res.ok) {
      deleteCookie("token")
      navigate("/v1/user/login")
    }
  }

  useEffect(() => {
    getUserInfo()
  }, [])

  return (
    <div className="wrapper">
      <div>
        {isFetching ? (
          <div>fetching details...</div>
        ) : (
          <div>
            {user && (
              <div>
                <h1>Welcome, {user && user.name}</h1>
                <p>{user && user.email}</p>
                <br />
                <button onClick={logout}>logout</button>
              </div>
            )}
          </div>
        )}

        <Errors errors={errors} />
      </div>
      <div className="card">
          <p>{data && (data.message)}</p>
          {loading && <div>A moment please...</div>}
          {error && (<div>{`There is a problem fetching the post data - ${error}`}</div>)}
        </div>
    </div>
  )
}

export default Session