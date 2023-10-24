const Summary = ({ data }) => {
  return (
    <ol>
      {data?.chain?.map((edge, idx) => {
        return (
          <li key={idx}>
            <TmdbPerson id={edge.src.id} name={edge.src.name} /> plays{" "}
            <em>{edge.src.Character}</em> in{" "}
            <TmdbMovie id={edge.src.Id} name={edge.src.Title} /> with{" "}
            <TmdbPerson id={edge.dest.id} name={edge.dest.name} /> who plays{" "}
            <em>{edge.dest.Character}</em>.
          </li>
        );
      })}
    </ol>
  );
};

const TmdbLink = ({ base, id, name }) => (
  <a target="_blank" rel="noreferrer" href={`${base}/${id}`}>
    {name}
  </a>
);

const TmdbPerson = (props) => (
  <TmdbLink base="https://www.themoviedb.org/person" {...props} />
);

const TmdbMovie = (props) => (
  <TmdbLink base="https://www.themoviedb.org/movie" {...props} />
);

export default Summary;
