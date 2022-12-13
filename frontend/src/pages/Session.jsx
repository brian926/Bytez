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
      if(success) {
        setUser(user)
      }
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

  const fakeUserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"
  
    const [post, setPost] = useState("") 
    const [responseToPost, setState] = useState(null) 
  
    const handleSubmit = async (e) => {
      e.preventDefault()
      try {
        const response = await fetch(Endpoints.createUrl, {
        method: 'POST',
        body: JSON.stringify({"long_url":post, "user_id": fakeUserId})
  
      })
      if (!response.ok) {
        throw new Error(
          `This is an HTTP error: The status is ${response.status}`
        );
      }
      let actualData = await response.json();
      console.log(actualData)
      setState(actualData.short_url)
      } catch(err) {
        console.log(err)
      }
    } 

  return (
    <div className="wrapper">
        <div>
          <a href="https://vitejs.dev" target="_blank">
            <img src="/vite.svg" className="logo" alt="Vite logo" />
          </a>
          <a href="https://reactjs.org" target="_blank">
            <img src="/react.svg" className="logo" alt="React logo" />
          </a>
          <a href="https://go.dev" target="_blank">
            <img src="/go.svg" className="logo" alt="Go logo" />
          </a>
          <a href="https://redis.io" target="_blank">
            <img src="/redis.svg" className="logo" alt="Redis logo" />
          </a>
          <a href="https://www.postgresql.org" target="_blank">
            <img src="/postgres.svg" className="logo" alt="Postgres logo" />
          </a>
        </div>
        <h1>Vite + React</h1>
        <p>Enter a URL link down below and generate a shorter cleaner URL for easy sharing and posting!</p>
      <div>
        {isFetching ? (
          <div>fetching details...</div>
        ) : (
          <div>
            {/* TODO: If there is user, pass userId and pull/save old urls */}
            {user && (
              <div>
                <div>
                  <h1>Welcome, {user && user.name}</h1>
                  <p>{user && user.email}</p>
                  <br />
                  <button onClick={logout}>logout</button>
                </div>
                <div>
                  <form onSubmit={handleSubmit}>
                    <input 
                      type="text"
                      value={post}
                      name="long_url"
                      onChange={ e => setPost(e.target.value)}
                    />
                    <button type="submit">Generate</button>
                  </form>
                  {responseToPost && (<div> <a target="_blank" href={responseToPost}>{responseToPost}</a></div>)}
                </div>
              </div>
            )}
            {/* Use fake userId to generate temp short urls */}
            <div>
                <form onSubmit={handleSubmit}>
                  <input 
                    type="text"
                    value={post}
                    name="long_url"
                    onChange={ e => setPost(e.target.value)}
                  />
                  <button type="submit">Generate</button>
                </form>
                {responseToPost && (<div> <a className="link" target="_blank" href={responseToPost}>{responseToPost}</a></div>)}
              </div>
          </div>
        )}

        <Errors errors={errors} />
      </div>
      <div className="card">
          <p>{data && (data.message)}</p>
          {loading && <div>A moment please...</div>}
          {error && (<div>{`There is a problem fetching the post data - ${error}`}</div>)}
      </div>
        <p className="read-the-docs">
          Click on the Vite, React, Go, Redis, and Postgres logos to learn more
        </p>
    </div>
  )
}

export default Session