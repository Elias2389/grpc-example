package database

import (
	"ae.com/proto-buffers/model"
	"context"
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
)

type PostgresRepository struct {
	db *sql.DB
}

const (
	host     = "localhost"
	port     = 5432
	user     = "root"
	password = "root"
	dbname   = "test-grpc-db"
)

func NewPostgresRepository() (*PostgresRepository, error) {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("You connected to your database.")

	return &PostgresRepository{db: db}, nil
}

func (repo *PostgresRepository) SetStudent(ctx context.Context, student *model.Student) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO students (id, name, age) VALUES ($1, $2, $3)",
		student.Id,
		student.Name,
		student.Age,
	)
	return err
}

func (repo *PostgresRepository) GetStudent(ctx context.Context, id string) (*model.Student, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, name, age FROM students WHERE id = $1",
		id,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var student = model.Student{}

	for rows.Next() {
		err := rows.Scan(&student.Id, &student.Name, &student.Age)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &student, nil
	}
	return &student, nil
}

func (repo *PostgresRepository) SetTest(ctx context.Context, test *model.Test) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO tests (id, name) VALUES ($1, $2)",
		test.Id,
		test.Name,
	)
	return err
}

func (repo *PostgresRepository) GetTest(ctx context.Context, id string) (*model.Test, error) {
	rows, err := repo.db.QueryContext(
		ctx,
		"SELECT id, name FROM tests WHERE id = $1",
		id,
	)

	if err != nil {
		return nil, err
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()

	var test = model.Test{}

	for rows.Next() {
		err := rows.Scan(&test.Id, &test.Name)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		return &test, nil
	}
	return &test, nil
}

func (repo *PostgresRepository) SetQuestion(ctx context.Context, test *model.Question) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO questions (id, answer, question, test_id) VALUES ($1, $2, $3, $4)",
		test.Id,
		test.Answer,
		test.Question,
		test.TestId,
	)
	return err
}

func (repo *PostgresRepository) SetEnrollment(ctx context.Context, enrollment *model.Enrollment) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO enrollments(student_id, test_id) VALUES($1, $2)", enrollment.StudentId, enrollment.TestId)
	return err
}

func (repo *PostgresRepository) GetStudentsPerTest(ctx context.Context, testId string) ([]*model.Student, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, name, age FROM students WHERE id IN (SELECT student_id FROM enrollments WHERE test_id = $1)", testId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var students []*model.Student
	for rows.Next() {
		var student = model.Student{}
		if err = rows.Scan(&student.Id, &student.Name, &student.Age); err == nil {
			students = append(students, &student)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return students, nil
}

func (repo *PostgresRepository) GetQuestionsPerTest(ctx context.Context, testId string) ([]*model.Question, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, question FROM questions WHERE test_id = $1", testId)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var questions []*model.Question
	for rows.Next() {
		var question = model.Question{}
		if err = rows.Scan(&question.Id, &question.Question); err == nil {
			questions = append(questions, &question)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return questions, nil
}
