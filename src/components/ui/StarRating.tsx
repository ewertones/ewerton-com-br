import { Star } from "lucide-react";
import styles from "./StarRating.module.css";

interface StarRatingProps {
    rating: number; // 0.0 to 5.0
    size?: number;
}

export function StarRating({ rating, size = 16 }: StarRatingProps) {
    // Ensure rating is between 0 and 5
    const clampedRating = Math.max(0, Math.min(5, rating));
    const percentage = (clampedRating / 5) * 100;

    const starsArray = [1, 2, 3, 4, 5];

    return (
        <div className={styles.starContainer} title={`${clampedRating.toFixed(1)} out of 5.0`}>
            {/* Background empty stars */}
            <div className={styles.starsBackground}>
                {starsArray.map((i) => (
                    <Star key={i} size={size} className={styles.starEmpty} />
                ))}
            </div>
            
            {/* Foreground filled stars with width mask */}
            <div className={styles.starsForeground} style={{ width: `${percentage}%` }}>
                {starsArray.map((i) => (
                    <Star key={i} size={size} className={styles.starFilled} fill="currentColor" />
                ))}
            </div>
        </div>
    );
}
