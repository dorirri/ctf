export default function ChallengeCard({ challenge, onClick }) {
  return (
    <div
      className={`challenge-card${challenge.is_solved ? ' solved' : ''}`}
      onClick={() => onClick(challenge)}
    >
      <div className="card-header">
        <span className="card-title">{challenge.title}</span>
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
    </div>
  )
}
