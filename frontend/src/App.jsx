import { useEffect, useState } from 'react'
import reactLogo from './assets/react.svg'
import './App.css'

function App() {

  // Get hello message
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


  // Create short url
  const UserId = "e0dba740-fc4b-4977-872c-d360239e6b1a"

  const [post, setPost] = useState("") 
  const [responseToPost, setState] = useState(null) 

  const handleSubmit = async (e) => {
    e.preventDefault()
    try {
      const response = await fetch(`http://localhost:9808/create-short-url`, {
      method: 'POST',
      body: JSON.stringify({"long_url":post, "user_id": UserId})

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
    <div className="App">
      <div>
        <a href="https://vitejs.dev" target="_blank">
          <img src="/vite.svg" className="logo" alt="Vite logo" />
        </a>
        <a href="https://reactjs.org" target="_blank">
          <img src={reactLogo} className="logo react" alt="React logo" />
        </a>
      </div>
      <h1>Vite + React</h1>
      <div className="card">
        <p>{data && (data.message)}</p>
        {loading && <div>A moment please...</div>}
        {error && (<div>{`There is a problem fetching the post data - ${error}`}</div>)}
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
      <p className="read-the-docs">
        Click on the Vite and React logos to learn more
      </p>
    </div>
  )
}

export default App
