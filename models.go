package main

type NbaPlayer struct {
	Object string `json:"object"`
	Id string `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Number string `json:"number" db:"number"`
	Team string `json:"team" db:"team"`
	Position string `json:"position" db:"pos"`
	Height int `json:"height" db:"height"`
	Weight int `json:"weight" db:"weight"`
}

type NbaGame struct {
	Object string `json:"object"`
	PlayerId string `json:"player_id" db:"player_id"`
	Date string `json:"date" db:"date"`
	Opponent string `json:"opponent" db:"opp"`
	Away int `json:"away" db:"away"`
	Score string `json:"score" db:"score"`
	SecondsPlayed int `json:"seconds_played" db:"sec_played"`
	FieldGoalsMade int `json:"field_goals_made" db:"fgm"`
	FieldGoalsAttempted int `json:"field_goals_attempted" db:"fga"`
	FieldGoalPercentage float32 `json:"field_goal_percentage" db:"fg_pct"`
	ThreePointersMade int `json:"three_pointers_made" db:"three_pm"`
	ThreePointersAttempted int `json:"three_pointers_attempted" db:"three_pa"`
	ThreePointPercentage float32 `json:"three_pointer_percentage" db:"three_pct"`
	FreeThrowsMade int `json:"free_throws_made" db:"ftm"`
	FreeThrowsAttempted int `json:"free_throws_attempted" db:"fta"`
	FreeThrowPercentage float32 `json:"free_throw_percentage" db:"ft_pct"`
	OffensiveRebounds int `json:"offensive_rebounds" db:"off_reb"`
	DefensiveRebounds int `json:"defensive_rebounds" db:"def_reb"`
	TotalRebounds int `json:"total_rebounds" db:"total_reb"`
	Assists int `json:"assists" db:"ast"`
	Turnovers int `json:"turnovers" db:"to"`
	Steals int `json:"steals" db:"stl"`
	Blocks int `json:"blocks" db:"blk"`
	PersonalFouls int `json:"personal_fouls" db:"pf"`
	Points int `json:"points" db:"pts"`
}

type NbaTeams struct {
	Object string `json:"object"`
	Teams []NbaTeam `json:"teams"`
}

type NbaTeam struct {
	Id string `json:"id" db:"id"`
}

type NbaRoster struct {
	Object string `json:"object"`
	Players []NbaPlayer `json:"players"`
}

type NbaGames struct {
	Object string `json:"object"`
	Games []NbaGame `json:"games"`
}

type NbaCategoryLeaders struct {
	Object string `json:"object"`
	Category string `json:"category"`
	Leaders []NbaCategoryLeader `json:"leaders"`
}

type NbaCategoryLeader struct {
	Id string `json:"id" db:"id"`
	CatAvg string `json:"value" db:"cat_avg"`
}

// NbaPlayer retrieves the personal data for a single NBA player
func (db *DB) NbaPlayer(player_id string) (NbaPlayer, error) {
	player := NbaPlayer{}
	player.Object = "nba_player"
	err := db.Get(&player, "SELECT id, name, number, team, pos, height, weight FROM nba_player WHERE id=?;", player_id)

	if err != nil {
		return player, err
	}

	return player, nil
}

// NbaCategoryLeaders retrieves the top 10 leaders in an NBA statistical category (single season)
func (db *DB) NbaCategoryLeaders(category string) (NbaCategoryLeaders, error) {
	categoryLeaders := NbaCategoryLeaders{"nba_categories", category, []NbaCategoryLeader{}}
	err := db.Select(&categoryLeaders.Leaders, "SELECT player_id AS id, AVG(" + category + ") AS cat_avg FROM nba_game GROUP BY player_id ORDER BY cat_avg DESC LIMIT 10")

	if err != nil {
		return categoryLeaders, err
	}

	return categoryLeaders, nil
}

// NbaTeams retrieves all active NBA teams in a single season
func (db *DB) NbaTeams() (NbaTeams, error) {
	teams := NbaTeams{"nba_teams", []NbaTeam{}}
	err := db.Select(&teams.Teams, "SELECT * FROM nba_team")

	if err != nil {
		return teams, err
	}

	return teams, nil
}

// NbaRoster retrieves all players on a single NBA team's roster
func (db *DB) NbaRoster(team_id string) (NbaRoster, error) {
	roster := NbaRoster{"nba_roster", []NbaPlayer{}}
	rows, err := db.Queryx("SELECT id, name, number, team, pos, height, weight FROM nba_player WHERE team=?;", team_id)

	if err != nil {
		return roster, err
	}

	for rows.Next() {
		var player NbaPlayer
		err = rows.StructScan(&player)
		player.Object = "nba_player"
		roster.Players = append(roster.Players, player)
	}

	if err = rows.Err(); err != nil {
		return roster, err
	}

	return roster, nil
}

// NbaGames retrieves all available game data for a single NBA player
func (db *DB) NbaGames(player_id string) (NbaGames, error) {
	games := NbaGames{"nba_games", []NbaGame{}}
	err := db.Select(&games.Games, "SELECT player_id, date, opp, away, COALESCE(score, '') as score, sec_played, fgm, fga, fg_pct, three_pm, three_pa, three_pct, ftm, fta, ft_pct, off_reb, def_reb, total_reb, ast, `to` FROM nba_game WHERE player_id=?;", player_id)

	if err != nil {
		return games, err
	}

	return games, nil
}
