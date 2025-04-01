package data

const (
	queryUserInsert = `
	INSERT INTO auth.users (name, surname, email, password_hash)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	queryUserGet = `
	SELECT id, name, surname, email, password_hash
	FROM auth.users
	WHERE id = $1
	`

	queryUserGetByEmail = `
	SELECT id, name, surname, email, password_hash
	FROM auth.users
	WHERE email = $1
	`

	queryProjectInsert = `
	INSERT INTO content.projects (title, access_code, creator_id)
	VALUES ($1, $2, $3)
	RETURNING id, access_code, created_at
	`

	queryProjectUserInsert = `
	INSERT INTO content.projects_users (project_id, user_id)
	VALUES ($1, $2)
	`

	queryProjectGetByAccessCode = `
	SELECT id, title, access_code, creator_id, created_at
	FROM content.projects
	WHERE access_code = $1
	`

	queryProjectGetAllForUser = `
	SELECT p.id, p.title, p.access_code, p.creator_id, p.created_at
	FROM content.projects p
	JOIN content.projects_users pu ON p.id = pu.project_id
	WHERE pu.user_id = $1;
	`

	queryProjectGet = `
	SELECT id, title, access_code, creator_id, created_at
	FROM content.projects
	WHERE id = $1
	`

	queryCourseInsert = `
	INSERT INTO content.courses (title, academic_year, project_id)
	VALUES ($1, $2, $3)
	RETURNING id
	`

	queryCourseGetAllForProject = `
	SELECT id, title, academic_year, project_id
	FROM content.courses
	WHERE project_id = $1;
	`

	queryCuratorInsert = `
	INSERT INTO content.curators (name, surname, birthday, status, project_id)
	VALUES ($1, $2, $3, $4)
	RETURNING id
	`

	queryCuratorsGetAllForProject = `
	SELECT id, name, surname, birthday, status, project_id
	FROM content.curators
	WHERE project_id = $1;
	`
)
