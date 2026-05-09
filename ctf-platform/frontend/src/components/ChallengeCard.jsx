export default function ChallengeCard({ challenge, onClick }) {
  return (
    <button
      type="button"
      className={`challenge-card${challenge.is_solved ? ' solved' : ''}`}
      onClick={() => onClick(challenge)}
      aria-label={`Open ${challenge.title || 'untitled challenge'}`}
    >
      <div className="card-header">
        <span className="card-title">{challenge.title || 'Untitled challenge'}</span>
        {challenge.is_solved && (
          <span className="card-solved-badge">SOLVED</span>
        )}
      </div>
      <div className="card-meta">
        <span className="card-category">{challenge.category}</span>
        <span>
          <span className="card-points">{challenge.points}</span>
          <span className="card-points-label"> pts</span>
        </span>
      </div>
    </button>
  )
}
