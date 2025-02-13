import { useEffect, useState } from "react";

function App() {
    const [movies, setMovies] = useState([]);

    useEffect(() => {
        console.log("Отправка запроса к API...");
    
        fetch("http://localhost:8080/movies")
            .then(response => response.json())
            .then(data => {
                console.log("Ответ от сервера:", data);
                setMovies(data.movies);
            })
            .catch(error => console.error("Ошибка запроса:", error));
    }, []);     

    console.log("Текущее состояние movies:", movies);

    return (
        <div className="container mt-5">
            <h1 className="text-center">Movies</h1>
            <div className="row mt-4">
                {movies.map(movie => (
                    <div key={movie.id} className="col-md-4 mb-4">
                        <div className="card shadow-sm">
                            <img src={movie.poster_url} className="card-img-top" alt={movie.title} />
                            <div className="card-body">
                                <h5 className="card-title">{movie.title} ({movie.year})</h5>
                                <p className="card-text"><strong>Genre:</strong> {movie.genre}</p>
                                <p className="card-text">{movie.description}</p>
                                <p className="card-text"><strong>Actors:</strong> {movie.actors.join(", ")}</p>
                            </div>
                        </div>
                    </div>
                ))}
            </div>
        </div>
    );
}

export default App;
