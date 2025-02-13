import React, { useEffect, useState } from "react";
import "bootstrap/dist/css/bootstrap.min.css";

function Movies() {
    const [movies, setMovies] = useState([]);

    useEffect(() => {
        fetch("http://localhost:8080/movies")
            .then(response => response.json())
            .then(data => {
                if (data.movies) {
                    setMovies(data.movies);
                }
            })
            .catch(error => console.error("Error fetching movies:", error));
    }, []);

    return (
        <div className="container mt-5">
            <h1 className="text-center">Movies</h1>
            <div className="row mt-4">
                {movies.map((movie, index) => (
                    <div key={index} className="col-md-4 mb-4">
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

export default Movies;
