import React, { useState, useEffect } from "react"
import { Endpoints } from "./api"
import Errors from "./components/Errors"
import './App.css'

function App() {
  const [isCopied, setIsCopied] = useState(false)
  async function copyTextToClipboard(text) {
    if ('clipboard' in navigator ) {
      return await navigator.clipboard.writeText(text)
    } else {
      return document.execCommand('copy', true, text)
    }
  }

  const handleCopyClick = () => {
    copyTextToClipboard(responseToPost)
    .then(() => {
      setIsCopied(true)
      setTimeout(() => {
        setIsCopied(false)
      }, 1500)
    })
    .catch((err) => {
      console.log(err)
    })
  }

  // Fetch welcome message, connecting to backend
  const [isFetching, setIsFetching] = useState(false)
  const [errors, setErrors] = useState([])

    const [data, setData] = useState(null)
    const [loading, setLoading] = useState(true)
    const [error, setError] = useState(null)

    useEffect(() => {
        const getData = async () => {
        try {
          setIsFetching(true)
          const response = await fetch(Endpoints.home);
          if (!response.ok) {
          throw new Error(
              `This is an HTTP error: The status is ${response.status}`
          );
          }
          let actualData = await response.json();
          setData(actualData);
          setError(null);
        } catch(err) {
            console.log(err)
            setError(err.message);
        } finally {
            setIsFetching(false)
            setLoading(false);
        }  
        }
        getData()
    }, [])

  // Handle URL creations
  const fakeUserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"
  
  const [post, setPost] = useState("") 
  const [responseToPost, setState] = useState(null) 

  const handleSubmit = async (e) => {
    e.preventDefault()
    if (post) {
      try {
        const response = await fetch(Endpoints.createUrl, {
        method: 'POST',
        body: JSON.stringify({"long_url":post, "user_id": fakeUserId})

        })
        let actualData = await response.json();
        if (!response.ok) {
          setErrors([actualData.error])
          throw new Error(
            `This is an HTTP error: The status is ${response.status}`
          );
        }
          setErrors([])
          setState(actualData.short_url)
      } catch(err) {
          setState(null)
          setErrors([actualData.error])
          console.log(err)
      }
    }
    else {
      setErrors(["Please enter a link"])
      setState(null)
    }
  } 

  return (
    <div className="wrapper">
        <div>
            <img src="/bytez-blue.png" className="bytez" alt="bytez" />
        </div>
        <div className="card">
          <p><b>{data && (data.message)}</b></p>
          {loading && <div>A moment please...</div>}
          {error && (<div>{`We ran into a problem - ${error}`}</div>)}
        </div>
        <p>Enter a URL link down below and generate a shorter cleaner URL for easy sharing and posting!</p>
      <div>
        {isFetching ? (
          <div>fetching details...</div>
        ) : (
          <div>
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
              </div>
            {responseToPost && (<div> <a className="link" target="_blank" href={responseToPost}>{responseToPost}</a><button onClick={handleCopyClick}><span>{isCopied ? 'Copied!' : 'Copy'}</span></button></div>)}
            <div><Errors errors={errors} /></div>
          </div>
        )}

        <div className="footer">
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
          <p className="read-the-docs">
            Click on the Vite, React, Go, Redis, and Postgres logos to learn more
          </p>
        </div>
      </div>
    </div>
  )
}

export default App
