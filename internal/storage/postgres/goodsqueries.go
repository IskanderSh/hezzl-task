package postgres

const createGoodQuery = `INSERT INTO goods (project_id, name, priority, removed) 
			VALUES ($1, $2, $3, $4) RETURNING *`

const getAllGoods = `SELECT * FROM goods`

const getGood = `SELECT * FROM goods WHERE id=$1 AND project_id=$2`

const updateGood = `UPDATE goods SET name=$1, description=$2 
             WHERE id=$3 AND project_id=$4 RETURNING *`

const deleteGood = `DELETE FROM goods WHERE id=$1 AND project_id=$2 RETURNING id, project_id, removed`

//const listGoods = `SELECT * FROM goods ORDER BY id LIMIT $1 OFFSET $2`
